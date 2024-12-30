package usecase_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/utils"
)

type FrecastsMessageReq struct {
	Area         float64 `json:"area"`
	HarvestTime  int     `json:"harvest_time"`
	HarvestYield float64 `json:"harvest_yield"`
	Demand       float64 `json:"demand"`
	Supply       float64 `json:"supply"`
	Sale         float64 `json:"sale"`
	Price        float64 `json:"price"`
	Day          int     `json:"day"`
}

type FrecastsMessageRes struct {
	OriginalData   FrecastsMessageReq `json:"original_data"`
	PredictedPrice float64            `json:"predicted_price"`
}

type ForecastsUsecaseImpl struct {
	landCommodityRepo repository_interface.LandCommodityRepository
	cityRepo          repository_interface.CityRepository
	priceRepo         repository_interface.PriceRepository
	priceHistoryRepo  repository_interface.PriceHistoryRepository
	demandRepo        repository_interface.DemandRepository
	demandHistoryRepo repository_interface.DemandHistoryRepository
	supplyRepo        repository_interface.SupplyRepository
	supplyHistoryRepo repository_interface.SupplyHistoryRepository
	saleRepo          repository_interface.SaleRepository
	harvestRepo       repository_interface.HarvestRepository
	commodityRepo     repository_interface.CommodityRepository
	rabbitMQ          messages.RabbitMQ
}

func NewForecastsUsecase(
	landCommodityRepo repository_interface.LandCommodityRepository,
	cityRepo repository_interface.CityRepository,
	priceRepo repository_interface.PriceRepository,
	priceHistoryRepo repository_interface.PriceHistoryRepository,
	demandRepo repository_interface.DemandRepository,
	demandHistoryRepo repository_interface.DemandHistoryRepository,
	supplyRepo repository_interface.SupplyRepository,
	supplyHistoryRepo repository_interface.SupplyHistoryRepository,
	saleRepo repository_interface.SaleRepository,
	harvestRepo repository_interface.HarvestRepository,
	commodityRepo repository_interface.CommodityRepository,
	rabbitMQ messages.RabbitMQ,
) usecase_interface.ForecastsUsecase {
	return &ForecastsUsecaseImpl{
		landCommodityRepo: landCommodityRepo,
		cityRepo:          cityRepo,
		priceRepo:         priceRepo,
		priceHistoryRepo:  priceHistoryRepo,
		demandRepo:        demandRepo,
		demandHistoryRepo: demandHistoryRepo,
		supplyRepo:        supplyRepo,
		supplyHistoryRepo: supplyHistoryRepo,
		saleRepo:          saleRepo,
		harvestRepo:       harvestRepo,
		commodityRepo:     commodityRepo,
		rabbitMQ:          rabbitMQ,
	}
}

func (f *ForecastsUsecaseImpl) GetForecastsByCommodityIDAndCityID(ctx context.Context, landCommodityID uuid.UUID, cityID int64) (*dto.ForecastsResponseDTO, error) {

	landCommodity, err := f.landCommodityRepo.FindByID(ctx, landCommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("land commodity not found")
	}

	commodityID := landCommodity.CommodityID

	commodity, err := f.commodityRepo.FindByID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	logrus.Log.Info("forecasts message sent", "commodity")

	city, err := f.cityRepo.FindByID(ctx, cityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}
	logrus.Log.Info("forecasts message sent", "city")

	demand, err := f.demandRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewNotFoundError("demand not found")
	}
	logrus.Log.Info("forecasts message sent", "demand")

	supply, err := f.supplyRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewNotFoundError("supply not found")
	}
	logrus.Log.Info("forecasts message sent", "supply")

	sales, err := f.saleRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewNotFoundError("sale not found")
	}
	logrus.Log.Info("forecasts message sent", sales)

	harvests, err := f.harvestRepo.FindByLandCommodityID(ctx, landCommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("harvest not found")
	}
	logrus.Log.Info("forecasts message sent", harvests)

	price, err := f.priceRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewNotFoundError("price not found")
	}
	logrus.Log.Info("forecasts message sent", price)

	var totalHarvest float64

	for _, h := range harvests {
		totalHarvest += h.Quantity
	}

	harvestYield := totalHarvest / float64(len(harvests))
	logrus.Log.Info("harvestYield", harvestYield)

	// Duration format: HH:MM:SS
	parts := strings.Split(commodity.Duration, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid duration format: expected HH:MM:SS, got %s", commodity.Duration)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid hours in duration: %w", err)
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minutes in duration: %w", err)
	}
	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid seconds in duration: %w", err)
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second

	harvestTime := landCommodity.CreatedAt.Add(duration)
	now := time.Now()
	daysUntilHarvest := int(harvestTime.Sub(now).Hours() / 24)

	logrus.Log.Info("forecasts message sent", "harvestTime")

	var sale float64

	for _, s := range sales {
		sale += s.Quantity
	}

	sale = sale / float64(len(sales))
	logrus.Log.Info("forecasts message sent", "saleQuantity")

	message := FrecastsMessageReq{
		Area:         landCommodity.LandArea,
		HarvestTime:  daysUntilHarvest,
		HarvestYield: harvestYield,
		Demand:       demand.Quantity,
		Supply:       supply.Quantity,
		Sale:         sale,
		Price:        price.Price,
		Day:          1,
	}

	logrus.Log.Info("forecasts message sent", message)

	err = f.rabbitMQ.PublishJSON(ctx, "prediction-exchange", "prediction-input", message)
	if err != nil {
		logrus.Log.Error(err)
		return nil, utils.NewInternalError(err.Error())
	}
	logrus.Log.Info("forecasts message sended")

	// Consume message from rabbitMQ
	msgs, err := f.rabbitMQ.ConsumeMessages("prediction-output-queue")
	if err != nil {
		logrus.Log.Error(err)
		return nil, utils.NewInternalError(err.Error())
	}

	// Tunggu pesan masuk
	var forecastsResponse dto.ForecastsResponseDTO

	var messageRes FrecastsMessageRes
	select {
	case msg := <-msgs:
		logrus.Log.Info("Message received: ", string(msg.Body))
		// res := json.Unmarshal(msg.Body, any)
		if err := json.Unmarshal(msg.Body, &messageRes); err != nil {
			logrus.Log.Error("Failed to unmarshal message: ", err)
			return nil, utils.NewInternalError("invalid message format")
		}
	case <-time.After(10 * time.Second): // Timeout 10 detik
		logrus.Log.Error("Timeout waiting for message from predict-output-queue")
		return nil, utils.NewInternalError("no response from prediction service")
	}

	forecastsResponse.HarvestPrice = messageRes.PredictedPrice
	forecastsResponse.HarvestDate = harvestTime
	forecastsResponse.City = city
	forecastsResponse.Commodity = commodity
	forecastsResponse.CurrentPrice = price.Price
	return &forecastsResponse, nil

}

func (f *ForecastsUsecaseImpl) GetForecastsByArea(ctx context.Context, area string) (*dto.ForecastsResponseDTO, error) {

	return nil, nil
}

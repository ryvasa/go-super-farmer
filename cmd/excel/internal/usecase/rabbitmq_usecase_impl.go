package usecase

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/cmd/excel/internal/repository"
)

type Message struct {
	Service string      `json:"Service"`
	Data    interface{} `json:"Data"`
}

type RabbitMQUsecaseImpl struct {
	reportRepo   repository.ReportRepository
	excelService ExcelInterface
}
type PriceMessage struct {
	CommodityID uuid.UUID `json:"CommodityID"`
	RegionID    uuid.UUID `json:"RegionID"`
	StartDate   time.Time `json:"StartDate"`
	EndDate     time.Time `json:"EndDate"`
}

type HarvestMessage struct {
	LandCommodityID uuid.UUID `json:"LandCommodityID"`
	StartDate       time.Time `json:"StartDate"`
	EndDate         time.Time `json:"EndDate"`
}

func NewRabbitMQUsecase(repo repository.ReportRepository, excelSvc ExcelInterface) RabbitMQUsecase {
	return &RabbitMQUsecaseImpl{
		reportRepo:   repo,
		excelService: excelSvc,
	}
}

func (u *RabbitMQUsecaseImpl) HandlePriceHistoryMessage(msgBody []byte) error {
	var msg PriceMessage
	if err := json.Unmarshal(msgBody, &msg); err != nil {
		return err
	}

	// Ambil data dari repository
	results, err := u.reportRepo.GetPriceHistoryReport(msg.StartDate, msg.EndDate, msg.CommodityID, msg.RegionID)
	if err != nil {
		return err
	}

	// Generate excel menggunakan usecase
	if err := u.excelService.CreatePriceHistoryReport(results, results[0].Commodity.Name, results[0].Region.City.Name); err != nil {
		return err
	}

	return nil
}

func (u *RabbitMQUsecaseImpl) HandleHarvestMessage(msgBody []byte) error {
	var msg HarvestMessage
	if err := json.Unmarshal(msgBody, &msg); err != nil {
		return err
	}
	log.Println("ahi")

	// Ambil data dari repository
	results, err := u.reportRepo.GetHarvestReport(msg.StartDate, msg.EndDate, msg.LandCommodityID)
	if err != nil {
		return err
	}

	// Generate excel menggunakan usecase
	if err := u.excelService.CreateHarvestReport(results, results[0].LandCommodity.Commodity.Name, results[0].Region.City.Name, results[0].LandCommodity.Land.User.Name); err != nil {
		return err
	}

	log.Println(results)

	return nil
}

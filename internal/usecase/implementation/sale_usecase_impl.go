package usecase_implementation

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/utils"
)

type SaleUsecaseImpl struct {
	saleRepo      repository_interface.SaleRepository
	cityRepo      repository_interface.CityRepository
	commodityRepo repository_interface.CommodityRepository
	cache         cache.Cache
}

func NewSaleUsecase(
	saleRepo repository_interface.SaleRepository,
	cityRepo repository_interface.CityRepository,
	commodityRepo repository_interface.CommodityRepository,
	cache cache.Cache,
) usecase_interface.SaleUsecase {
	return &SaleUsecaseImpl{
		saleRepo:      saleRepo,
		cityRepo:      cityRepo,
		commodityRepo: commodityRepo,
		cache:         cache,
	}
}

func (uc *SaleUsecaseImpl) CreateSale(ctx context.Context, req *dto.SaleCreateDTO) (*domain.Sale, error) {
	sale := domain.Sale{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	_, err := uc.cityRepo.FindByID(ctx, req.CityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}
	_, err = uc.commodityRepo.FindByID(ctx, req.CommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	logrus.Log.Info(req)
	parseDate, err := time.Parse("2006-01-02", req.SaleDate)
	if err != nil {
		return nil, utils.NewBadRequestError("sale date format is invalid")
	}

	sale.CityID = req.CityID
	sale.CommodityID = req.CommodityID
	sale.Quantity = req.Quantity
	sale.Unit = req.Unit
	sale.Price = req.Price
	sale.SaleDate = parseDate
	sale.ID = uuid.New()

	err = uc.saleRepo.Create(ctx, &sale)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdSale, err := uc.saleRepo.FindByID(ctx, sale.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdSale, nil
}

func (uc *SaleUsecaseImpl) GetAllSales(ctx context.Context, params *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
	if err := params.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	sales, err := uc.saleRepo.FindAll(ctx, params)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	count, err := uc.saleRepo.Count(ctx, &params.Filter)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	response := &dto.PaginationResponseDTO{
		TotalRows:  count,
		TotalPages: int(math.Ceil(float64(count) / float64(params.Limit))),
		Page:       params.Page,
		Limit:      params.Limit,
		Data:       sales,
	}
	return response, nil
}

func (uc *SaleUsecaseImpl) GetSaleByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	sale, err := uc.saleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return sale, nil
}

func (uc *SaleUsecaseImpl) GetSalesByCommodityID(ctx context.Context, params *dto.PaginationDTO, id uuid.UUID) (*dto.PaginationResponseDTO, error) {
	if err := params.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	params.Filter.CommodityID = &id

	sales, err := uc.saleRepo.FindByCommodityID(ctx, params, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	count, err := uc.saleRepo.Count(ctx, &params.Filter)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	response := &dto.PaginationResponseDTO{
		TotalRows:  count,
		TotalPages: int(math.Ceil(float64(count) / float64(params.Limit))),
		Page:       params.Page,
		Limit:      params.Limit,
		Data:       sales,
	}
	return response, nil
}

func (uc *SaleUsecaseImpl) GetSalesByCityID(ctx context.Context, params *dto.PaginationDTO, id int64) (*dto.PaginationResponseDTO, error) {
	if err := params.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	params.Filter.CityID = &id

	sales, err := uc.saleRepo.FindByCityID(ctx, params, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	count, err := uc.saleRepo.Count(ctx, &params.Filter)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	response := &dto.PaginationResponseDTO{
		TotalRows:  count,
		TotalPages: int(math.Ceil(float64(count) / float64(params.Limit))),
		Page:       params.Page,
		Limit:      params.Limit,
		Data:       sales,
	}
	return response, nil
}

func (uc *SaleUsecaseImpl) UpdateSale(ctx context.Context, id uuid.UUID, req *dto.SaleUpdateDTO) (*domain.Sale, error) {
	sale := &domain.Sale{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	sale, err := uc.saleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("sale not found")
	}

	if req.SaleDate != "" {
		parseDate, err := time.Parse("2006-01-02", req.SaleDate)
		if err != nil {
			return nil, utils.NewBadRequestError("sale date format is invalid")
		}
		sale.SaleDate = parseDate
	}

	sale.CommodityID = req.CommodityID
	sale.Quantity = req.Quantity
	sale.Unit = req.Unit
	sale.Price = req.Price
	sale.CityID = req.CityID

	err = uc.saleRepo.Update(ctx, id, sale)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedSale, err := uc.saleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedSale, nil
}

func (uc *SaleUsecaseImpl) DeleteSale(ctx context.Context, id uuid.UUID) error {
	_, err := uc.saleRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	err = uc.saleRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (uc *SaleUsecaseImpl) RestoreSale(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	_, err := uc.saleRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	err = uc.saleRepo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	sale, err := uc.saleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return sale, nil
}

func (uc *SaleUsecaseImpl) GetAllDeletedSales(ctx context.Context, params *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
	if err := params.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	sales, err := uc.saleRepo.FindAllDeleted(ctx, params)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	count, err := uc.saleRepo.Count(ctx, &params.Filter)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	response := &dto.PaginationResponseDTO{
		TotalRows:  count,
		TotalPages: int(math.Ceil(float64(count) / float64(params.Limit))),
		Page:       params.Page,
		Limit:      params.Limit,
		Data:       sales,
	}
	return response, nil
}

func (uc *SaleUsecaseImpl) GetDeletedSaleByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	sale, err := uc.saleRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return sale, nil
}

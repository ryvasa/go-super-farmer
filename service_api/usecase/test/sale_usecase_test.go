package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/service_api/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/service_api/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type SaleRepoMock struct {
	Sale      *mock_repo.MockSaleRepository
	City      *mock_repo.MockCityRepository
	Commodity *mock_repo.MockCommodityRepository
	Cache     *mock_pkg.MockCache
}

type SaleIDs struct {
	SaleID      uuid.UUID
	CommodityID uuid.UUID
	CityID      int64
}

type SaleDomains struct {
	Sale      *domain.Sale
	Sales     []*domain.Sale
	City      *domain.City
	Commodity *domain.Commodity
}

type SaleDtos struct {
	Create             *dto.SaleCreateDTO
	Update             *dto.SaleUpdateDTO
	Pagination         *dto.PaginationDTO
	Filter             *dto.ParamFilterDTO
	PaginationResponse *dto.PaginationResponseDTO
}

func SaleUsecaseSetup(t *testing.T) (*SaleIDs, *SaleDomains, *SaleDtos, *SaleRepoMock, usecase_interface.SaleUsecase, context.Context) {
	cityID := int64(1)
	saleID := uuid.New()
	commodityID := uuid.New()
	ids := &SaleIDs{
		SaleID:      saleID,
		CommodityID: commodityID,
		CityID:      cityID,
	}

	date := "2022-01-01"
	saleDate, _ := time.Parse("2006-01-02", date)
	domains := &SaleDomains{
		Sale: &domain.Sale{
			ID:          saleID,
			CommodityID: commodityID,
			CityID:      cityID,
			Quantity:    float64(100),
			Unit:        "kg",
			Price:       float64(10),
			SaleDate:    saleDate,
		},
		City: &domain.City{
			ID:   cityID,
			Name: "kota",
		},
		Sales: []*domain.Sale{
			{
				ID:          saleID,
				CommodityID: commodityID,
				CityID:      cityID,
				Quantity:    float64(100),
				Unit:        "kg",
				Price:       float64(10),
				SaleDate:    saleDate,
			},
		},
		Commodity: &domain.Commodity{
			ID:   commodityID,
			Name: "komodit",
		},
	}
	dtos := &SaleDtos{
		Create: &dto.SaleCreateDTO{
			CommodityID: commodityID,
			CityID:      cityID,
			Quantity:    float64(100),
			Unit:        "kg",
			Price:       float64(10),
			SaleDate:    date,
		},
		Update: &dto.SaleUpdateDTO{
			Quantity: float64(200),
			Unit:     "kg",
			Price:    float64(20),
			SaleDate: date,
		},
		Filter: &dto.ParamFilterDTO{
			CommodityID: &ids.CommodityID,
			CityID:      &ids.CityID,
		},
		Pagination: &dto.PaginationDTO{
			Page:  1,
			Limit: 10,
			Filter: dto.ParamFilterDTO{
				CommodityID: &ids.CommodityID,
				CityID:      &ids.CityID,
			},
		},
		PaginationResponse: &dto.PaginationResponseDTO{
			TotalRows:  1,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       domains.Sales,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock_repo.NewMockCityRepository(ctrl)
	commodityRepo := mock_repo.NewMockCommodityRepository(ctrl)
	saleRepo := mock_repo.NewMockSaleRepository(ctrl)
	cache := mock_pkg.NewMockCache(ctrl)

	usecase := usecase_implementation.NewSaleUsecase(saleRepo, cityRepo, commodityRepo, cache)
	ctx := context.Background()

	repo := &SaleRepoMock{
		Sale:      saleRepo,
		City:      cityRepo,
		Commodity: commodityRepo,
		Cache:     cache,
	}

	return ids, domains, dtos, repo, usecase, ctx
}

func TestSaleUsecase_CreateSale(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should create sale successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Sale.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Sale) error {
			p.ID = ids.SaleID
			return nil
		}).Times(1)

		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		resp, err := uc.CreateSale(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.CommodityID, resp.CommodityID)
		assert.Equal(t, dtos.Create.CityID, resp.CityID)
		assert.Equal(t, dtos.Create.Quantity, resp.Quantity)
		assert.Equal(t, dtos.Create.Unit, resp.Unit)
		assert.Equal(t, dtos.Create.Price, resp.Price)
		parsed, err := time.Parse("2006-01-02", dtos.Create.SaleDate)
		assert.NoError(t, err)
		assert.Equal(t, parsed, resp.SaleDate)
	})

	t.Run("should return error if validation fails", func(t *testing.T) {
		resp, err := uc.CreateSale(ctx, &dto.SaleCreateDTO{
			CommodityID: ids.CommodityID,
			CityID:      ids.CityID,
			Quantity:    -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error if commodity not found", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.CreateSale(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error if city not found", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.CreateSale(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error if date is invalid", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		resp, err := uc.CreateSale(ctx, &dto.SaleCreateDTO{
			CommodityID: ids.CommodityID,
			CityID:      ids.CityID,
			Quantity:    100,
			Unit:        "kg",
			Price:       10,
			SaleDate:    "122",
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "sale date format is invalid")
	})

	t.Run("should return error if create fails", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Sale.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateSale(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error if find by id fails", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Sale.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Sale) error {
			p.ID = ids.SaleID
			return nil
		}).Times(1)

		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateSale(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_GetAllSales(t *testing.T) {
	_, domains, dtos, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should get all sales successfully", func(t *testing.T) {

		repo.Sale.EXPECT().FindAll(ctx, gomock.Any()).Return(domains.Sales, nil).Times(1)

		repo.Sale.EXPECT().Count(ctx, gomock.Any()).Return(int64(len(domains.Sales)), nil).Times(1)

		resp, err := uc.GetAllSales(ctx, dtos.Pagination)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Pagination.Page, resp.Page)
		assert.Equal(t, dtos.Pagination.Limit, resp.Limit)
		assert.Equal(t, dtos.PaginationResponse.TotalRows, resp.TotalRows)
		assert.Equal(t, dtos.PaginationResponse.TotalPages, resp.TotalPages)
		assert.Equal(t, dtos.PaginationResponse.Data, resp.Data)
	})
	t.Run("should return validation error", func(t *testing.T) {
		invalidQueryParams := &dto.PaginationDTO{
			Limit: -1, // Invalid limit
			Page:  1,
		}

		resp, err := uc.GetAllSales(ctx, invalidQueryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "limit must be greater than 0")
	})

	t.Run("should return error if find all fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindAll(ctx, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllSales(ctx, dtos.Pagination)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_GetSaleByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should get sale by id successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		resp, err := uc.GetSaleByID(ctx, ids.SaleID)

		assert.NoError(t, err)
		assert.Equal(t, domains.Sale, resp)
	})

	t.Run("should return error if find by id fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSaleByID(ctx, ids.SaleID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_GetSalesByCommodityID(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should get sales by commodity id successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindByCommodityID(ctx, gomock.Any(), ids.CommodityID).Return(domains.Sales, nil).Times(1)

		repo.Sale.EXPECT().Count(ctx, gomock.Any()).Return(int64(len(domains.Sales)), nil).Times(1)

		resp, err := uc.GetSalesByCommodityID(ctx, dtos.Pagination, ids.CommodityID)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Pagination.Page, resp.Page)
		assert.Equal(t, dtos.Pagination.Limit, resp.Limit)
		assert.Equal(t, dtos.PaginationResponse.TotalRows, resp.TotalRows)
		assert.Equal(t, dtos.PaginationResponse.TotalPages, resp.TotalPages)
		assert.Equal(t, dtos.PaginationResponse.Data, resp.Data)
	})

	t.Run("should return error if find by commodity id fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindByCommodityID(ctx, gomock.Any(), ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSalesByCommodityID(ctx, dtos.Pagination, ids.CommodityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_GetSalesByCityID(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should get sales by city id successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindByCityID(ctx, gomock.Any(), ids.CityID).Return(domains.Sales, nil).Times(1)

		repo.Sale.EXPECT().Count(ctx, gomock.Any()).Return(int64(len(domains.Sales)), nil).Times(1)

		resp, err := uc.GetSalesByCityID(ctx, dtos.Pagination, ids.CityID)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Pagination.Page, resp.Page)
		assert.Equal(t, dtos.Pagination.Limit, resp.Limit)
		assert.Equal(t, dtos.PaginationResponse.TotalRows, resp.TotalRows)
		assert.Equal(t, dtos.PaginationResponse.TotalPages, resp.TotalPages)
		assert.Equal(t, dtos.PaginationResponse.Data, resp.Data)
	})

	t.Run("should return error if find by city id fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindByCityID(ctx, gomock.Any(), ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSalesByCityID(ctx, dtos.Pagination, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_UpdateSale(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should update sale successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		repo.Sale.EXPECT().Update(ctx, ids.SaleID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Sale) error {
			p.ID = ids.SaleID
			return nil
		}).Times(1)

		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		resp, err := uc.UpdateSale(ctx, ids.SaleID, dtos.Update)

		assert.NoError(t, err)
		assert.Equal(t, domains.Sale, resp)
	})

	t.Run("should return error if find by id fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(nil, utils.NewNotFoundError("sale not found")).Times(1)

		resp, err := uc.UpdateSale(ctx, ids.SaleID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "sale not found")
	})

	t.Run("should return error if update fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		repo.Sale.EXPECT().Update(ctx, ids.SaleID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateSale(ctx, ids.SaleID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_DeleteSale(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should delete sale successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		repo.Sale.EXPECT().Delete(ctx, ids.SaleID).Return(nil).Times(1)

		err := uc.DeleteSale(ctx, ids.SaleID)

		assert.NoError(t, err)
	})

	t.Run("should return error if find by id fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteSale(ctx, ids.SaleID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error if delete fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		repo.Sale.EXPECT().Delete(ctx, ids.SaleID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteSale(ctx, ids.SaleID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_RestoreSale(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should restore sale successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindDeletedByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		repo.Sale.EXPECT().Restore(ctx, ids.SaleID).Return(nil).Times(1)

		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		resp, err := uc.RestoreSale(ctx, ids.SaleID)

		assert.NoError(t, err)
		assert.Equal(t, domains.Sale, resp)
		assert.Equal(t, ids.SaleID, resp.ID)
		assert.Equal(t, ids.CommodityID, resp.CommodityID)
		assert.Equal(t, ids.CityID, resp.CityID)
		assert.Equal(t, domains.Sale.Quantity, resp.Quantity)
	})

	t.Run("should return error if find deleted by id not found", func(t *testing.T) {
		repo.Sale.EXPECT().FindDeletedByID(ctx, ids.SaleID).Return(nil, utils.NewNotFoundError("sale not found")).Times(1)

		resp, err := uc.RestoreSale(ctx, ids.SaleID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "sale not found")
	})

	t.Run("should return error if restore fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindDeletedByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		repo.Sale.EXPECT().Restore(ctx, ids.SaleID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreSale(ctx, ids.SaleID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error if find restored by id fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindDeletedByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		repo.Sale.EXPECT().Restore(ctx, ids.SaleID).Return(nil).Times(1)

		repo.Sale.EXPECT().FindByID(ctx, ids.SaleID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreSale(ctx, ids.SaleID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_GetAllDeletedSales(t *testing.T) {
	_, domains, dtos, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should get all deleted sales successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindAllDeleted(ctx, gomock.Any()).Return(domains.Sales, nil).Times(1)

		repo.Sale.EXPECT().Count(ctx, gomock.Any()).Return(int64(len(domains.Sales)), nil).Times(1)

		resp, err := uc.GetAllDeletedSales(ctx, dtos.Pagination)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Pagination.Page, resp.Page)
		assert.Equal(t, dtos.Pagination.Limit, resp.Limit)
		assert.Equal(t, dtos.PaginationResponse.TotalRows, resp.TotalRows)
		assert.Equal(t, dtos.PaginationResponse.TotalPages, resp.TotalPages)
		assert.Equal(t, dtos.PaginationResponse.Data, resp.Data)
	})

	t.Run("should return validation error", func(t *testing.T) {
		invalidQueryParams := &dto.PaginationDTO{
			Limit: -1, // Invalid limit
			Page:  1,
		}

		resp, err := uc.GetAllDeletedSales(ctx, invalidQueryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "limit must be greater than 0")
	})

	t.Run("should return error if find all fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindAllDeleted(ctx, gomock.Any()).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllDeletedSales(ctx, dtos.Pagination)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSaleUsecase_GetDeletedSaleByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SaleUsecaseSetup(t)

	t.Run("should get deleted sale by id successfully", func(t *testing.T) {
		repo.Sale.EXPECT().FindDeletedByID(ctx, ids.SaleID).Return(domains.Sale, nil).Times(1)

		resp, err := uc.GetDeletedSaleByID(ctx, ids.SaleID)

		assert.NoError(t, err)
		assert.Equal(t, domains.Sale, resp)
	})

	t.Run("should return error if find deleted by id fails", func(t *testing.T) {
		repo.Sale.EXPECT().FindDeletedByID(ctx, ids.SaleID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetDeletedSaleByID(ctx, ids.SaleID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

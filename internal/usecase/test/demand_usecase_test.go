package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type DemandRepoMock struct {
	Demand        *mock_repo.MockDemandRepository
	DemandHistory *mock_repo.MockDemandHistoryRepository
	Commodity     *mock_repo.MockCommodityRepository
	City          *mock_repo.MockCityRepository
	TxManager     *mock_pkg.MockTransactionManager
}

type DemandIDs struct {
	DemandID        uuid.UUID
	CommodityID     uuid.UUID
	CityID          int64
	DemandHistoryID uuid.UUID
}

type DemandDomainMocks struct {
	Demand         *domain.Demand
	Demands        []*domain.Demand
	UpdatedDemand  *domain.Demand
	Commodity      *domain.Commodity
	City           *domain.City
	DemandHistorys []*domain.DemandHistory
}

type DemandDTOMocks struct {
	Create *dto.DemandCreateDTO
	Update *dto.DemandUpdateDTO
}

func DemandUsecaseSetup(t *testing.T) (*DemandIDs, *DemandDomainMocks, *DemandDTOMocks, *DemandRepoMock, usecase_interface.DemandUsecase, context.Context) {
	cityID := int64(1)
	commodityID := uuid.New()
	demandID := uuid.New()
	demandHstoryID := uuid.New()

	ids := &DemandIDs{
		DemandID:        demandID,
		CommodityID:     commodityID,
		CityID:          cityID,
		DemandHistoryID: demandHstoryID,
	}

	domains := &DemandDomainMocks{
		Demand: &domain.Demand{
			ID:          demandID,
			CommodityID: commodityID,
			CityID:      cityID,
			Quantity:    10,
		},
		Demands: []*domain.Demand{
			{
				ID:          demandID,
				CommodityID: commodityID,
				CityID:      cityID,
				Quantity:    10,
			},
		},
		UpdatedDemand: &domain.Demand{
			ID:          demandID,
			CommodityID: commodityID,
			CityID:      cityID,
			Quantity:    20,
		},
		Commodity: &domain.Commodity{
			ID: commodityID,
		},
		City: &domain.City{
			ID: cityID,
		},
		DemandHistorys: []*domain.DemandHistory{
			{
				ID:          demandHstoryID,
				CommodityID: commodityID,
				CityID:      cityID,
				Quantity:    20,
			},
		},
	}

	dtos := &DemandDTOMocks{
		Create: &dto.DemandCreateDTO{
			CommodityID: commodityID,
			CityID:      cityID,
			Quantity:    10,
		},
		Update: &dto.DemandUpdateDTO{
			Quantity: 20,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock_repo.NewMockCityRepository(ctrl)
	commodityRepo := mock_repo.NewMockCommodityRepository(ctrl)
	demandRepo := mock_repo.NewMockDemandRepository(ctrl)
	demandHistoryRepo := mock_repo.NewMockDemandHistoryRepository(ctrl)
	txRepo := mock_pkg.NewMockTransactionManager(ctrl)

	uc := usecase_implementation.NewDemandUsecase(demandRepo, demandHistoryRepo, commodityRepo, cityRepo, txRepo)
	ctx := context.Background()

	repo := &DemandRepoMock{
		Demand:        demandRepo,
		City:          cityRepo,
		Commodity:     commodityRepo,
		DemandHistory: demandHistoryRepo,
		TxManager:     txRepo,
	}

	return ids, domains, dtos, repo, uc, ctx
}

func TestDemandUsecase_CreateDemand(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should create demand successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Demand.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Demand) error {
			p.ID = ids.DemandID
			return nil
		}).Times(1)

		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, ids.CommodityID, resp.CommodityID)
		assert.Equal(t, ids.CityID, resp.CityID)
		assert.Equal(t, ids.DemandID, resp.ID)
	})

	t.Run("should return error when validation fails", func(t *testing.T) {

		resp, err := uc.CreateDemand(ctx, &dto.DemandCreateDTO{
			CommodityID: ids.CommodityID,
			CityID:      ids.CityID,
			Quantity:    -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when city not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when create demand fails", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Demand.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandUsecase_GetAllDemands(t *testing.T) {
	_, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should get all demands successfully", func(t *testing.T) {
		repo.Demand.EXPECT().FindAll(ctx).Return(domains.Demands, nil).Times(1)

		resp, err := uc.GetAllDemands(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Demands), len(resp))
		assert.Equal(t, (domains.Demands)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get all demands fails", func(t *testing.T) {
		repo.Demand.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllDemands(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandUsecase_GetDemandByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should get demand by id successfully", func(t *testing.T) {
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		resp, err := uc.GetDemandByID(ctx, ids.DemandID)

		assert.NoError(t, err)
		assert.Equal(t, ids.DemandID, resp.ID)
	})

	t.Run("should return error when get demand by id fails", func(t *testing.T) {
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetDemandByID(ctx, ids.DemandID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandUsecase_GetDemandsByCommodityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should get demands by commodity id successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Demand.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(domains.Demands, nil).Times(1)

		resp, err := uc.GetDemandsByCommodityID(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Demands), len(resp))
		assert.Equal(t, (domains.Demands)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get demands by commodity id fails", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Demand.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetDemandsByCommodityID(ctx, ids.CommodityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandUsecase_GetDemandsByCityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should get demands by city id successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Demand.EXPECT().FindByCityID(ctx, ids.CityID).Return(domains.Demands, nil).Times(1)

		resp, err := uc.GetDemandsByCityID(ctx, ids.CityID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Demands), len(resp))
		assert.Equal(t, (domains.Demands)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get demands by city id fails", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Demand.EXPECT().FindByCityID(ctx, ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetDemandsByCityID(ctx, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when city not found", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.GetDemandsByCityID(ctx, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})
}

func TestDemandUsecase_UpdateDemand(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should update demand successfully", func(t *testing.T) {
		// Setup mock untuk WithTransaction
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		// Setup mock untuk operasi dalam transaction
		repo.Demand.EXPECT().
			FindByID(ctx, ids.DemandID).
			Return(domains.Demand, nil)

		repo.DemandHistory.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, history *domain.DemandHistory) error {
				history.ID = ids.DemandHistoryID
				return nil
			})

		repo.Demand.EXPECT().
			Update(ctx, ids.DemandID, gomock.Any()).
			Return(nil)

		repo.Demand.EXPECT().
			FindByID(ctx, ids.DemandID).
			Return(domains.UpdatedDemand, nil)

		// Execute
		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.DemandID, resp.ID)
		assert.Equal(t, float64(20), resp.Quantity)
	})

	t.Run("should rollback transaction when create history fails", func(t *testing.T) {
		// Setup mock untuk WithTransaction
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				err := fn(ctx)
				return err
			})

		// Setup mock untuk operasi dalam transaction
		repo.Demand.EXPECT().
			FindByID(ctx, ids.DemandID).
			Return(domains.Demand, nil)

		repo.DemandHistory.EXPECT().
			Create(ctx, gomock.Any()).
			Return(fmt.Errorf("create history error"))

		// Execute
		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "create history error")
	})

	t.Run("should rollback transaction when update demand fails", func(t *testing.T) {
		// Setup mock untuk WithTransaction
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				err := fn(ctx)
				return err
			})

		// Setup mock untuk operasi dalam transaction
		repo.Demand.EXPECT().
			FindByID(ctx, ids.DemandID).
			Return(domains.Demand, nil)

		repo.DemandHistory.EXPECT().
			Create(ctx, gomock.Any()).
			Return(nil)

		repo.Demand.EXPECT().
			Update(ctx, ids.DemandID, gomock.Any()).
			Return(fmt.Errorf("update demand error"))

		// Execute
		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "update demand error")
	})

	t.Run("should return error when validation fails", func(t *testing.T) {

		resp, err := uc.UpdateDemand(ctx, ids.DemandID, &dto.DemandUpdateDTO{
			Quantity: -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when demand not found", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(nil, utils.NewNotFoundError("demand not found")).Times(1)

		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "demand not found")
	})

	t.Run("should return error when get updated demand fails", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		repo.DemandHistory.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

		repo.Demand.EXPECT().Update(ctx, ids.DemandID, gomock.Any()).Return(nil).Times(1)

		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandUsecase_DeleteDemand(t *testing.T) {
	ids, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should delete demand successfully", func(t *testing.T) {
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		repo.Demand.EXPECT().Delete(ctx, ids.DemandID).Return(nil).Times(1)

		err := uc.DeleteDemand(ctx, ids.DemandID)

		assert.NoError(t, err)
	})

	t.Run("should return error when demand not found", func(t *testing.T) {
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(nil, utils.NewNotFoundError("demand not found")).Times(1)

		err := uc.DeleteDemand(ctx, ids.DemandID)

		assert.Error(t, err)
		assert.EqualError(t, err, "demand not found")
	})

	t.Run("should return error when delete demand fails", func(t *testing.T) {
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		repo.Demand.EXPECT().Delete(ctx, ids.DemandID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteDemand(ctx, ids.DemandID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandUsecase_GetDemandHistoryByCommodityIDAndCityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should return demand history successfully", func(t *testing.T) {

		repo.DemandHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(domains.DemandHistorys, nil).Times(1)

		repo.Demand.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(domains.Demand, nil).Times(1)

		resp, err := uc.GetDemandHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp))
		assert.Equal(t, (domains.DemandHistorys)[0].ID, (resp)[0].ID)
		assert.Equal(t, (domains.DemandHistorys)[0].Quantity, (resp)[0].Quantity)
	})

	t.Run("should return error when get demand history fails", func(t *testing.T) {
		repo.DemandHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetDemandHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when demand not found", func(t *testing.T) {
		repo.DemandHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(domains.DemandHistorys, nil).Times(1)

		repo.Demand.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetDemandHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

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

type SupplyRepoMock struct {
	Supply        *mock_repo.MockSupplyRepository
	SupplyHistory *mock_repo.MockSupplyHistoryRepository
	Commodity     *mock_repo.MockCommodityRepository
	City          *mock_repo.MockCityRepository
	TxManager     *mock_pkg.MockTransactionManager
}

type SupplyIDs struct {
	SupplyID        uuid.UUID
	CommodityID     uuid.UUID
	CityID          int64
	SupplyHistoryID uuid.UUID
}

type SupplyDomainMocks struct {
	Supply         *domain.Supply
	Supplys        []*domain.Supply
	UpdatedSupply  *domain.Supply
	Commodity      *domain.Commodity
	City           *domain.City
	SupplyHistorys []*domain.SupplyHistory
}

type SupplyDTOMocks struct {
	Create *dto.SupplyCreateDTO
	Update *dto.SupplyUpdateDTO
}

func SupplyUsecaseSetup(t *testing.T) (*SupplyIDs, *SupplyDomainMocks, *SupplyDTOMocks, *SupplyRepoMock, usecase_interface.SupplyUsecase, context.Context) {
	cityID := int64(1)
	commodityID := uuid.New()
	supplyID := uuid.New()
	supplyHstoryID := uuid.New()

	ids := &SupplyIDs{
		SupplyID:        supplyID,
		CommodityID:     commodityID,
		CityID:          cityID,
		SupplyHistoryID: supplyHstoryID,
	}

	domains := &SupplyDomainMocks{
		Supply: &domain.Supply{
			ID:          supplyID,
			CommodityID: commodityID,
			CityID:      cityID,
			Quantity:    10,
		},
		Supplys: []*domain.Supply{
			{
				ID:          supplyID,
				CommodityID: commodityID,
				CityID:      cityID,
				Quantity:    10,
			},
		},
		UpdatedSupply: &domain.Supply{
			ID:          supplyID,
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
		SupplyHistorys: []*domain.SupplyHistory{
			{
				ID:          supplyHstoryID,
				CommodityID: commodityID,
				CityID:      cityID,
				Quantity:    50,
			},
		},
	}

	dtos := &SupplyDTOMocks{
		Create: &dto.SupplyCreateDTO{
			CommodityID: commodityID,
			CityID:      cityID,
			Quantity:    10,
		},
		Update: &dto.SupplyUpdateDTO{
			Quantity: 20,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cityRepo := mock_repo.NewMockCityRepository(ctrl)
	commodityRepo := mock_repo.NewMockCommodityRepository(ctrl)
	supplyRepo := mock_repo.NewMockSupplyRepository(ctrl)
	supplyHistoryRepo := mock_repo.NewMockSupplyHistoryRepository(ctrl)
	txRepo := mock_pkg.NewMockTransactionManager(ctrl)

	uc := usecase_implementation.NewSupplyUsecase(supplyRepo, supplyHistoryRepo, commodityRepo, cityRepo, txRepo)
	ctx := context.Background()

	repo := &SupplyRepoMock{
		Supply:        supplyRepo,
		City:          cityRepo,
		Commodity:     commodityRepo,
		SupplyHistory: supplyHistoryRepo,
		TxManager:     txRepo,
	}

	return ids, domains, dtos, repo, uc, ctx
}

func TestSupplyRepository_CreateSupply(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should create supply successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Supply.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Supply) error {
			p.ID = ids.SupplyID
			return nil
		}).Times(1)

		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(domains.Supply, nil).Times(1)

		resp, err := uc.CreateSupply(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, ids.CommodityID, resp.CommodityID)
		assert.Equal(t, ids.CityID, resp.CityID)
		assert.Equal(t, ids.SupplyID, resp.ID)
	})

	t.Run("should return error when validation fails", func(t *testing.T) {

		resp, err := uc.CreateSupply(ctx, &dto.SupplyCreateDTO{
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

		resp, err := uc.CreateSupply(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when city not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(nil, utils.NewNotFoundError("city not found")).Times(1)

		resp, err := uc.CreateSupply(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "city not found")
	})

	t.Run("should return error when create supply fails", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Supply.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateSupply(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSupplyRepository_GetAllSupply(t *testing.T) {
	_, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should get all supplys successfully", func(t *testing.T) {
		repo.Supply.EXPECT().FindAll(ctx).Return(domains.Supplys, nil).Times(1)

		resp, err := uc.GetAllSupply(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Supplys), len(resp))
		assert.Equal(t, (domains.Supplys)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get all supplys fails", func(t *testing.T) {
		repo.Supply.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllSupply(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSupplyRepository_GetSupplyByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should get supply by id successfully", func(t *testing.T) {
		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(domains.Supply, nil).Times(1)

		resp, err := uc.GetSupplyByID(ctx, ids.SupplyID)

		assert.NoError(t, err)
		assert.Equal(t, ids.SupplyID, resp.ID)
	})

	t.Run("should return error when get supply by id fails", func(t *testing.T) {
		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSupplyByID(ctx, ids.SupplyID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSupplyRepository_GetSupplyByCommodityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should get supplys by commodity id successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Supply.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(domains.Supplys, nil).Times(1)

		resp, err := uc.GetSupplyByCommodityID(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Supplys), len(resp))
		assert.Equal(t, (domains.Supplys)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get supplys by commodity id fails", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Supply.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSupplyByCommodityID(ctx, ids.CommodityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSupplyRepository_GetSupplyByCityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should get supplys by city id successfully", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Supply.EXPECT().FindByCityID(ctx, ids.CityID).Return(domains.Supplys, nil).Times(1)

		resp, err := uc.GetSupplyByCityID(ctx, ids.CityID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Supplys), len(resp))
		assert.Equal(t, (domains.Supplys)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get supplys by city id fails", func(t *testing.T) {
		repo.City.EXPECT().FindByID(ctx, ids.CityID).Return(domains.City, nil).Times(1)

		repo.Supply.EXPECT().FindByCityID(ctx, ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSupplyByCityID(ctx, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSupplyRepository_UpdateSupply(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should update supply successfully", func(t *testing.T) {
		// Setup mock untuk WithTransaction
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		// Setup mock untuk operasi dalam transaction
		repo.Supply.EXPECT().
			FindByID(ctx, ids.SupplyID).
			Return(domains.Supply, nil)

		repo.SupplyHistory.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, history *domain.SupplyHistory) error {
				history.ID = ids.SupplyHistoryID
				return nil
			})

		repo.Supply.EXPECT().
			Update(ctx, ids.SupplyID, gomock.Any()).
			Return(nil)

		repo.Supply.EXPECT().
			FindByID(ctx, ids.SupplyID).
			Return(domains.UpdatedSupply, nil)

		// Execute
		resp, err := uc.UpdateSupply(ctx, ids.SupplyID, dtos.Update)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.SupplyID, resp.ID)
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
		repo.Supply.EXPECT().
			FindByID(ctx, ids.SupplyID).
			Return(domains.Supply, nil)

		repo.SupplyHistory.EXPECT().
			Create(ctx, gomock.Any()).
			Return(fmt.Errorf("create history error"))

		// Execute
		resp, err := uc.UpdateSupply(ctx, ids.SupplyID, dtos.Update)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "create history error")
	})

	t.Run("should rollback transaction when update supply fails", func(t *testing.T) {
		// Setup mock untuk WithTransaction
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				err := fn(ctx)
				return err
			})

		// Setup mock untuk operasi dalam transaction
		repo.Supply.EXPECT().
			FindByID(ctx, ids.SupplyID).
			Return(domains.Supply, nil)

		repo.SupplyHistory.EXPECT().
			Create(ctx, gomock.Any()).
			Return(nil)

		repo.Supply.EXPECT().
			Update(ctx, ids.SupplyID, gomock.Any()).
			Return(fmt.Errorf("update supply error"))

		// Execute
		resp, err := uc.UpdateSupply(ctx, ids.SupplyID, dtos.Update)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "update supply error")
	})

	t.Run("should return error when validation fails", func(t *testing.T) {

		resp, err := uc.UpdateSupply(ctx, ids.SupplyID, &dto.SupplyUpdateDTO{
			Quantity: -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when supply not found", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(nil, utils.NewNotFoundError("supply not found")).Times(1)

		resp, err := uc.UpdateSupply(ctx, ids.SupplyID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "supply not found")
	})

	t.Run("should return error when get updated supply fails", func(t *testing.T) {
		repo.TxManager.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(domains.Supply, nil).Times(1)

		repo.SupplyHistory.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(1)

		repo.Supply.EXPECT().Update(ctx, ids.SupplyID, gomock.Any()).Return(nil).Times(1)

		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateSupply(ctx, ids.SupplyID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSupplyRepository_DeleteSupply(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should delete supply successfully", func(t *testing.T) {
		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(domains.Supply, nil).Times(1)

		repo.Supply.EXPECT().Delete(ctx, ids.SupplyID).Return(nil).Times(1)

		err := uc.DeleteSupply(ctx, ids.SupplyID)

		assert.NoError(t, err)
	})

	t.Run("should return error when supply not found", func(t *testing.T) {
		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(nil, utils.NewNotFoundError("supply not found")).Times(1)

		err := uc.DeleteSupply(ctx, ids.SupplyID)

		assert.Error(t, err)
		assert.EqualError(t, err, "supply not found")
	})

	t.Run("should return error when delete supply fails", func(t *testing.T) {
		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(domains.Supply, nil).Times(1)

		repo.Supply.EXPECT().Delete(ctx, ids.SupplyID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteSupply(ctx, ids.SupplyID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestSupplyUsecase_GetSupplyHistoryByCommodityIDAndCityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should return supply history successfully", func(t *testing.T) {

		repo.SupplyHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(domains.SupplyHistorys, nil).Times(1)

		repo.Supply.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(domains.Supply, nil).Times(1)

		resp, err := uc.GetSupplyHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp))
		assert.Equal(t, (domains.SupplyHistorys)[0].ID, (resp)[0].ID)
		assert.Equal(t, (domains.SupplyHistorys)[0].Quantity, (resp)[0].Quantity)
	})

	t.Run("should return error when get supply history fails", func(t *testing.T) {
		repo.SupplyHistory.EXPECT().FindByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSupplyHistoryByCommodityIDAndCityID(ctx, ids.CommodityID, ids.CityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

}

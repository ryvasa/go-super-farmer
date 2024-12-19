package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type SupplyRepoMock struct {
	Supply        *mock.MockSupplyRepository
	SupplyHistory *mock.MockSupplyHistoryRepository
	Commodity     *mock.MockCommodityRepository
	Region        *mock.MockRegionRepository
	TxManager     *mock_pkg.MockTransactionManager
}

type SupplyIDs struct {
	SupplyID        uuid.UUID
	CommodityID     uuid.UUID
	RegionID        uuid.UUID
	SupplyHistoryID uuid.UUID
}

type SupplyDomainMocks struct {
	Supply         *domain.Supply
	Supplys        []*domain.Supply
	UpdatedSupply  *domain.Supply
	Commodity      *domain.Commodity
	Region         *domain.Region
	SupplyHistorys []*domain.SupplyHistory
}

type SupplyDTOMocks struct {
	Create *dto.SupplyCreateDTO
	Update *dto.SupplyUpdateDTO
}

func SupplyUsecaseSetup(t *testing.T) (*SupplyIDs, *SupplyDomainMocks, *SupplyDTOMocks, *SupplyRepoMock, usecase_interface.SupplyUsecase, context.Context) {
	regionID := uuid.New()
	commodityID := uuid.New()
	supplyID := uuid.New()
	supplyHstoryID := uuid.New()

	ids := &SupplyIDs{
		SupplyID:        supplyID,
		CommodityID:     commodityID,
		RegionID:        regionID,
		SupplyHistoryID: supplyHstoryID,
	}

	domains := &SupplyDomainMocks{
		Supply: &domain.Supply{
			ID:          supplyID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Quantity:    10,
		},
		Supplys: []*domain.Supply{
			{
				ID:          supplyID,
				CommodityID: commodityID,
				RegionID:    regionID,
				Quantity:    10,
			},
		},
		UpdatedSupply: &domain.Supply{
			ID:          supplyID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Quantity:    20,
		},
		Commodity: &domain.Commodity{
			ID: commodityID,
		},
		Region: &domain.Region{
			ID: regionID,
		},
		SupplyHistorys: []*domain.SupplyHistory{
			{
				ID:          supplyHstoryID,
				CommodityID: commodityID,
				RegionID:    regionID,
				Quantity:    50,
			},
		},
	}

	dtos := &SupplyDTOMocks{
		Create: &dto.SupplyCreateDTO{
			CommodityID: commodityID,
			RegionID:    regionID,
			Quantity:    10,
		},
		Update: &dto.SupplyUpdateDTO{
			Quantity: 20,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	supplyRepo := mock.NewMockSupplyRepository(ctrl)
	supplyHistoryRepo := mock.NewMockSupplyHistoryRepository(ctrl)
	txRepo := mock_pkg.NewMockTransactionManager(ctrl)

	uc := usecase_implementation.NewSupplyUsecase(supplyRepo, supplyHistoryRepo, commodityRepo, regionRepo, txRepo)
	ctx := context.Background()

	repo := &SupplyRepoMock{
		Supply:        supplyRepo,
		Region:        regionRepo,
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

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Supply.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Supply) error {
			p.ID = ids.SupplyID
			return nil
		}).Times(1)

		repo.Supply.EXPECT().FindByID(ctx, ids.SupplyID).Return(domains.Supply, nil).Times(1)

		resp, err := uc.CreateSupply(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, ids.CommodityID, resp.CommodityID)
		assert.Equal(t, ids.RegionID, resp.RegionID)
		assert.Equal(t, ids.SupplyID, resp.ID)
	})

	t.Run("should return error when validation fails", func(t *testing.T) {

		resp, err := uc.CreateSupply(ctx, &dto.SupplyCreateDTO{
			CommodityID: ids.CommodityID,
			RegionID:    ids.RegionID,
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

	t.Run("should return error when region not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(nil, utils.NewNotFoundError("region not found")).Times(1)

		resp, err := uc.CreateSupply(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "region not found")
	})

	t.Run("should return error when create supply fails", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

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

func TestSupplyRepository_GetSupplyByRegionID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should get supplys by region id successfully", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Supply.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(domains.Supplys, nil).Times(1)

		resp, err := uc.GetSupplyByRegionID(ctx, ids.RegionID)

		assert.NoError(t, err)
		assert.Equal(t, len(domains.Supplys), len(resp))
		assert.Equal(t, (domains.Supplys)[0].ID, (resp)[0].ID)
	})

	t.Run("should return error when get supplys by region id fails", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Supply.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSupplyByRegionID(ctx, ids.RegionID)

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

func TestSupplyUsecase_GetSupplyHistoryByCommodityIDAndRegionID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := SupplyUsecaseSetup(t)

	t.Run("should return supply history successfully", func(t *testing.T) {

		repo.SupplyHistory.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(domains.SupplyHistorys, nil).Times(1)

		repo.Supply.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(domains.Supply, nil).Times(1)

		resp, err := uc.GetSupplyHistoryByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp))
		assert.Equal(t, (domains.SupplyHistorys)[0].ID, (resp)[0].ID)
		assert.Equal(t, (domains.SupplyHistorys)[0].Quantity, (resp)[0].Quantity)
	})

	t.Run("should return error when get supply history fails", func(t *testing.T) {
		repo.SupplyHistory.EXPECT().FindByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetSupplyHistoryByCommodityIDAndRegionID(ctx, ids.CommodityID, ids.RegionID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

}

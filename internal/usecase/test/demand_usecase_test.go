package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type DemandRepoMock struct {
	Demand        *mock.MockDemandRepository
	DemandHistory *mock.MockDemandHistoryRepository
	Commodity     *mock.MockCommodityRepository
	Region        *mock.MockRegionRepository
}

type DemandIDs struct {
	DemandID        uuid.UUID
	CommodityID     uuid.UUID
	RegionID        uuid.UUID
	DemandHistoryID uuid.UUID
}

type DemandDomainMocks struct {
	Demand        *domain.Demand
	Demands       *[]domain.Demand
	UpdatedDemand *domain.Demand
	Commodity     *domain.Commodity
	Region        *domain.Region
}

type DemandDTOMocks struct {
	Create *dto.DemandCreateDTO
	Update *dto.DemandUpdateDTO
}

func DemandUsecaseSetup(t *testing.T) (*DemandIDs, *DemandDomainMocks, *DemandDTOMocks, *DemandRepoMock, usecase.DemandUsecase, context.Context) {
	regionID := uuid.New()
	commodityID := uuid.New()
	demandID := uuid.New()
	demandHstoryID := uuid.New()

	ids := &DemandIDs{
		DemandID:        demandID,
		CommodityID:     commodityID,
		RegionID:        regionID,
		DemandHistoryID: demandHstoryID,
	}

	domains := &DemandDomainMocks{
		Demand: &domain.Demand{
			ID:          demandID,
			CommodityID: commodityID,
			RegionID:    regionID,
			Quantity:    10,
		},
		Demands: &[]domain.Demand{
			{
				ID:          demandID,
				CommodityID: commodityID,
				RegionID:    regionID,
				Quantity:    10,
			},
		},
		UpdatedDemand: &domain.Demand{
			ID:          demandID,
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
	}

	dtos := &DemandDTOMocks{
		Create: &dto.DemandCreateDTO{
			CommodityID: commodityID,
			RegionID:    regionID,
			Quantity:    10,
		},
		Update: &dto.DemandUpdateDTO{
			Quantity: 20,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	demandRepo := mock.NewMockDemandRepository(ctrl)
	demandHistoryRepo := mock.NewMockDemandHistoryRepository(ctrl)

	uc := usecase.NewDemandUsecase(demandRepo, demandHistoryRepo, commodityRepo, regionRepo)
	ctx := context.Background()

	repo := &DemandRepoMock{
		Demand:        demandRepo,
		Region:        regionRepo,
		Commodity:     commodityRepo,
		DemandHistory: demandHistoryRepo,
	}

	return ids, domains, dtos, repo, uc, ctx
}

func TestDemandRepository_CreateDemand(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should create demand successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Demand.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Demand) error {
			p.ID = ids.DemandID
			return nil
		}).Times(1)

		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, ids.CommodityID, resp.CommodityID)
		assert.Equal(t, ids.RegionID, resp.RegionID)
		assert.Equal(t, ids.DemandID, resp.ID)
	})

	t.Run("should return error when validation fails", func(t *testing.T) {

		resp, err := uc.CreateDemand(ctx, &dto.DemandCreateDTO{
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

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when region not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(nil, utils.NewNotFoundError("region not found")).Times(1)

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "region not found")
	})

	t.Run("should return error when create demand fails", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Demand.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateDemand(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandRepository_GetAllDemands(t *testing.T) {
	_, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should get all demands successfully", func(t *testing.T) {
		repo.Demand.EXPECT().FindAll(ctx).Return(domains.Demands, nil).Times(1)

		resp, err := uc.GetAllDemands(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Demands), len(*resp))
		assert.Equal(t, (*domains.Demands)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get all demands fails", func(t *testing.T) {
		repo.Demand.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllDemands(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandRepository_GetDemandByID(t *testing.T) {
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

func TestDemandRepository_GetDemandsByCommodityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should get demands by commodity id successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Demand.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(domains.Demands, nil).Times(1)

		resp, err := uc.GetDemandsByCommodityID(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Demands), len(*resp))
		assert.Equal(t, (*domains.Demands)[0].ID, (*resp)[0].ID)
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

func TestDemandRepository_GetDemandsByRegionID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should get demands by region id successfully", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Demand.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(domains.Demands, nil).Times(1)

		resp, err := uc.GetDemandsByRegionID(ctx, ids.RegionID)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Demands), len(*resp))
		assert.Equal(t, (*domains.Demands)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get demands by region id fails", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Demand.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetDemandsByRegionID(ctx, ids.RegionID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandRepository_UpdateDemand(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := DemandUsecaseSetup(t)

	t.Run("should update demand successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(domains.Commodity, nil).Times(1)

		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		repo.DemandHistory.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.DemandHistory) error {
			p.ID = ids.DemandHistoryID
			return nil
		}).Times(1)

		repo.Demand.EXPECT().Update(ctx, ids.DemandID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Demand) error {
			p.Quantity = float64(20)
			return nil
		}).Times(1)

		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.UpdatedDemand, nil).Times(1)

		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		assert.NoError(t, err)
		assert.Equal(t, ids.DemandID, resp.ID)
		assert.Equal(t, float64(20), resp.Quantity)
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
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(nil, utils.NewNotFoundError("demand not found")).Times(1)

		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "demand not found")
	})

	t.Run("should return error when update demand fails", func(t *testing.T) {
		repo.Demand.EXPECT().FindByID(ctx, ids.DemandID).Return(domains.Demand, nil).Times(1)

		repo.DemandHistory.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateDemand(ctx, ids.DemandID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDemandRepository_DeleteDemand(t *testing.T) {
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

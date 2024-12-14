package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type HarvestRepoMock struct {
	Harvest       *mock.MockHarvestRepository
	Region        *mock.MockRegionRepository
	LandCommodity *mock.MockLandCommodityRepository
}

type HarvestIDs struct {
	HarvestID       uuid.UUID
	LandCommodityID uuid.UUID
	RegionID        uuid.UUID
	LandID          uuid.UUID
	CommodityID     uuid.UUID
}

type HarvestDomainMock struct {
	Harvest        *domain.Harvest
	Harvests       *[]domain.Harvest
	UpdatedHarvest *domain.Harvest
	Region         *domain.Region
	LandCommodity  *domain.LandCommodity
}

type HarvestDTOMock struct {
	Create *dto.HarvestCreateDTO
	Update *dto.HarvestUpdateDTO
}

func HarvestUsecaseSetup(t *testing.T) (*HarvestIDs, *HarvestDomainMock, *HarvestDTOMock, *HarvestRepoMock, usecase.HarvestUsecase, context.Context) {
	regionID := uuid.New()
	landCommodityID := uuid.New()
	harvestID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	ids := &HarvestIDs{
		HarvestID:       harvestID,
		LandCommodityID: landCommodityID,
		RegionID:        regionID,
		LandID:          landID,
		CommodityID:     commodityID,
	}
	date, _ := time.Parse("2006-01-02", "2022-01-01")

	domains := &HarvestDomainMock{
		Harvest: &domain.Harvest{
			ID:              harvestID,
			LandCommodityID: landCommodityID,
			RegionID:        regionID,
			Quantity:        float64(100),
			Unit:            "kg",
			HarvestDate:     date,
		},
		Harvests: &[]domain.Harvest{
			{
				ID:              harvestID,
				LandCommodityID: landCommodityID,
				RegionID:        regionID,
				Quantity:        float64(100),
				Unit:            "kg",
				HarvestDate:     date,
			},
		},
		UpdatedHarvest: &domain.Harvest{
			ID:              harvestID,
			LandCommodityID: landCommodityID,
			RegionID:        regionID,
			Quantity:        float64(99),
			Unit:            "kg",
			HarvestDate:     date,
		},
		Region: &domain.Region{
			ID: regionID,
		},
		LandCommodity: &domain.LandCommodity{
			ID: landCommodityID,
		},
	}

	dto := &HarvestDTOMock{
		Create: &dto.HarvestCreateDTO{
			LandCommodityID: landCommodityID,
			RegionID:        regionID,
			Quantity:        float64(100),
			Unit:            "kg",
			HarvestDate:     "2022-01-01",
		},
		Update: &dto.HarvestUpdateDTO{
			Quantity: 99,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	regionRepo := mock.NewMockRegionRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	harvestRepo := mock.NewMockHarvestRepository(ctrl)

	uc := usecase.NewHarvestUsecase(harvestRepo, regionRepo, landCommodityRepo)
	ctx := context.TODO()

	repo := &HarvestRepoMock{Harvest: harvestRepo, Region: regionRepo, LandCommodity: landCommodityRepo}

	return ids, domains, dto, repo, uc, ctx
}

func TestHarvestRepository_CreateHarvest(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should create harvest successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Harvest.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
		assert.Equal(t, ids.LandCommodityID, resp.LandCommodityID)
		assert.Equal(t, ids.RegionID, resp.RegionID)
		assert.Equal(t, float64(100), resp.Quantity)
		assert.Equal(t, "kg", resp.Unit)
		assert.Equal(t, domains.Harvest.HarvestDate, resp.HarvestDate)
		assert.Equal(t, domains.Harvest.Quantity, resp.Quantity)
		assert.Equal(t, domains.Harvest.Unit, resp.Unit)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.CreateHarvest(ctx, &dto.HarvestCreateDTO{
			LandCommodityID: ids.LandCommodityID,
			RegionID:        ids.RegionID,
			Quantity:        -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when land commodity not found", func(t *testing.T) {

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("should return error when region not found", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(nil, utils.NewNotFoundError("region not found")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "region not found")
	})

	t.Run("should return error when create harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Harvest.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get created harvest ", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Harvest.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateHarvest(ctx, dtos.Create)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetAllHarvest(t *testing.T) {
	_, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get all harvests successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindAll(ctx).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetAllHarvest(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Harvests), len(*resp))
		assert.Equal(t, (*domains.Harvests)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get all harvests fails", func(t *testing.T) {
		repo.Harvest.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllHarvest(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvest by id successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		resp, err := uc.GetHarvestByID(ctx, ids.HarvestID)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
	})

	t.Run("should return error when get harvest by id not found", func(t *testing.T) {
		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		resp, err := uc.GetHarvestByID(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "harvest not found")
	})
}

// TODO: fix this test
func TestHarvestUsecase_GetHarvestByCommodityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by commodity id successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByCommodityID(ctx, ids.LandCommodityID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByCommodityID(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Harvests), len(*resp))
		assert.Equal(t, (*domains.Harvests)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get harvests by commodity id fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByCommodityID(ctx, ids.LandCommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByCommodityID(ctx, ids.LandCommodityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

// TODO: fix this test
func TestHarvestUsecase_GetHarvestByLandID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by land id successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByCommodityID(ctx, ids.LandCommodityID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByCommodityID(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Harvests), len(*resp))
		assert.Equal(t, (*domains.Harvests)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get harvests by land id fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByLandID(ctx, ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByLandID(ctx, ids.LandID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestByLandCommodityID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by land commodity id successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByLandCommodityID(ctx, ids.LandCommodityID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByLandCommodityID(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Harvests), len(*resp))
		assert.Equal(t, (*domains.Harvests)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get harvests by land commodity id fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByLandCommodityID(ctx, ids.LandCommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByLandCommodityID(ctx, ids.LandCommodityID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestByRegionID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get harvests by region id successfully", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Harvest.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetHarvestByRegionID(ctx, ids.RegionID)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Harvests), len(*resp))
		assert.Equal(t, (*domains.Harvests)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get harvests by region id fails", func(t *testing.T) {
		repo.Region.EXPECT().FindByID(ctx, ids.RegionID).Return(domains.Region, nil).Times(1)

		repo.Harvest.EXPECT().FindByRegionID(ctx, ids.RegionID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestByRegionID(ctx, ids.RegionID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_UpdateHarvest(t *testing.T) {
	ids, domains, dtos, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should update harvest successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.UpdatedHarvest, nil).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
		assert.Equal(t, float64(99), resp.Quantity)
		assert.Equal(t, "kg", resp.Unit)
	})

	t.Run("should return error when harvest not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "harvest not found")
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, &dto.HarvestUpdateDTO{
			Quantity: -10,
		})

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when update harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get updated harvest error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Update(ctx, ids.HarvestID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, p *domain.Harvest) error {
			p.ID = ids.HarvestID
			return nil
		}).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateHarvest(ctx, ids.HarvestID, dtos.Update)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_DeleteHarvest(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should delete harvest successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Delete(ctx, ids.HarvestID).Return(nil).Times(1)

		err := uc.DeleteHarvest(ctx, ids.HarvestID)

		assert.NoError(t, err)
	})

	t.Run("should return error when harvest not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("harvest not found")).Times(1)

		err := uc.DeleteHarvest(ctx, ids.HarvestID)

		assert.Error(t, err)
		assert.EqualError(t, err, "harvest not found")
	})

	t.Run("should return error when delete harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Delete(ctx, ids.HarvestID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteHarvest(ctx, ids.HarvestID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_RestoreHarvest(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should restore harvest successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(nil).Times(1)

		repo.Harvest.EXPECT().FindByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
	})

	t.Run("should return error when harvest not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(nil, utils.NewNotFoundError("deleted harvest not found")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "deleted harvest not found")
	})

	t.Run("should return error when restore harvest fails", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get restored harvest", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(domains.LandCommodity, nil).Times(1)

		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		repo.Harvest.EXPECT().Restore(ctx, ids.HarvestID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreHarvest(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetAllDeletedHarvest(t *testing.T) {
	_, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get all deleted harvests successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindAllDeleted(ctx).Return(domains.Harvests, nil).Times(1)

		resp, err := uc.GetAllDeletedHarvest(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(*domains.Harvests), len(*resp))
		assert.Equal(t, (*domains.Harvests)[0].ID, (*resp)[0].ID)
	})

	t.Run("should return error when get all deleted harvests fails", func(t *testing.T) {
		repo.Harvest.EXPECT().FindAllDeleted(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllDeletedHarvest(ctx)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestHarvestUsecase_GetHarvestDeletedByID(t *testing.T) {
	ids, domains, _, repo, uc, ctx := HarvestUsecaseSetup(t)

	t.Run("should get deleted harvest by id successfully", func(t *testing.T) {
		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(domains.Harvest, nil).Times(1)

		resp, err := uc.GetHarvestDeletedByID(ctx, ids.HarvestID)

		assert.NoError(t, err)
		assert.Equal(t, ids.HarvestID, resp.ID)
	})

	t.Run("should return error when get deleted harvest by id fails", func(t *testing.T) {
		repo.Harvest.EXPECT().FindDeletedByID(ctx, ids.HarvestID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetHarvestDeletedByID(ctx, ids.HarvestID)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.EqualError(t, err, "deleted harvest not found")
	})
}

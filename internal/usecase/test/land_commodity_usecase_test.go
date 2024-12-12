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

type LandCommodityRepoMock struct {
	LandCommodity *mock.MockLandCommodityRepository
	Land          *mock.MockLandRepository
	Commodity     *mock.MockCommodityRepository
}

type LandCommodityIDs struct {
	CommodityID     uuid.UUID
	LandID          uuid.UUID
	LandCommodityID uuid.UUID
}

type LandCommodityMocks struct {
	LandCommodity        *domain.LandCommodity
	LandCommodities      *[]domain.LandCommodity
	UpdatedLandCommodity *domain.LandCommodity
	Land                 *domain.Land
	Commodity            *domain.Commodity
}

type LandCommodityDTOMocks struct {
	Create *dto.LandCommodityCreateDTO
	Update *dto.LandCommodityUpdateDTO
}

func LandCommodityUtils(t *testing.T) (*LandCommodityIDs, *LandCommodityMocks, *LandCommodityDTOMocks, *LandCommodityRepoMock, usecase.LandCommodityUsecase, context.Context) {
	landID := uuid.New()
	commodityID := uuid.New()
	landCommodityID := uuid.New()

	ids := &LandCommodityIDs{
		CommodityID:     commodityID,
		LandID:          landID,
		LandCommodityID: landCommodityID,
	}

	mocks := &LandCommodityMocks{
		LandCommodity: &domain.LandCommodity{
			ID:          landCommodityID,
			CommodityID: commodityID,
			LandID:      landID,
			LandArea:    float64(100),
		},
		LandCommodities: &[]domain.LandCommodity{
			{
				ID:          landCommodityID,
				CommodityID: commodityID,
				LandID:      landID,
				LandArea:    float64(100),
			},
			{
				ID:          landCommodityID,
				CommodityID: commodityID,
				// LandID:      landID,
				//LandArea:    float64(100),
			},
		},
		UpdatedLandCommodity: &domain.LandCommodity{
			ID:          landCommodityID,
			CommodityID: commodityID,
			LandID:      landID,
			LandArea:    float64(200),
		},
		Land: &domain.Land{
			ID:       landID,
			LandArea: float64(1000),
		},
		Commodity: &domain.Commodity{
			ID: commodityID,
		},
	}

	dtoMocks := &LandCommodityDTOMocks{
		Create: &dto.LandCommodityCreateDTO{
			LandID:      landID,
			CommodityID: commodityID,
			LandArea:    float64(100),
		},
		Update: &dto.LandCommodityUpdateDTO{
			LandID:      landID,
			CommodityID: commodityID,
			LandArea:    float64(200),
		},
	}

	repoMock := &LandCommodityRepoMock{
		LandCommodity: mock.NewMockLandCommodityRepository(gomock.NewController(t)),
		Land:          mock.NewMockLandRepository(gomock.NewController(t)),
		Commodity:     mock.NewMockCommodityRepository(gomock.NewController(t)),
	}

	uc := usecase.NewLandCommodityUsecase(repoMock.LandCommodity, repoMock.Land, repoMock.Commodity)
	ctx := context.Background()

	return ids, mocks, dtoMocks, repoMock, uc, ctx
}
func TestLandCommodityUsecase_CreateLandCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := LandCommodityUtils(t)
	t.Run("should create land commodity successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(10), nil).Times(1)

		repo.LandCommodity.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, l *domain.LandCommodity) error {
			l.ID = ids.LandCommodityID
			return nil
		}).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
	})

	t.Run("should return error when land area is greater than max area", func(t *testing.T) {

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(1000), nil).Times(1)

		req := &dto.LandCommodityCreateDTO{LandID: ids.LandID, CommodityID: ids.CommodityID, LandArea: float64(100)}
		resp, err := uc.CreateLandCommodity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("should return error validation error", func(t *testing.T) {
		resp, err := uc.CreateLandCommodity(ctx, &dto.LandCommodityCreateDTO{LandID: ids.LandID, CommodityID: ids.CommodityID, LandArea: 0})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})
}

func TestLandCommodityUsecase_GetLandCommodityByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return error when land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		resp, err := uc.GetLandCommodityByID(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land commodity not found")
	})
	t.Run("should return land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		resp, err := uc.GetLandCommodityByID(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
		assert.Equal(t, ids.CommodityID, resp.CommodityID)
		assert.Equal(t, ids.LandID, resp.LandID)
		assert.Equal(t, float64(100), resp.LandArea)
	})
}

func TestLandCommodityUsecase_GetLandCommodityByLandID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByLandID(ctx, ids.LandID).Return(mocks.LandCommodities, nil).Times(1)

		resp, err := uc.GetLandCommodityByLandID(ctx, ids.LandID)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, (*mocks.LandCommodities)[0].LandArea, (*resp)[0].LandArea)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByLandID(ctx, ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetLandCommodityByLandID(ctx, ids.LandID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestLandCommodityUsecase_GetLandCommodityByCommodityID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return land commodity successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(mocks.LandCommodities, nil).Times(1)

		resp, err := uc.GetLandCommodityByCommodityID(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, (*mocks.LandCommodities)[0].LandArea, (*resp)[0].LandArea)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByCommodityID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetLandCommodityByCommodityID(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestLandCommodityUsecase_GetAllLandCommodity(t *testing.T) {
	_, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should return land commodity successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindAll(ctx).Return(mocks.LandCommodities, nil).Times(1)

		resp, err := uc.GetAllLandCommodity(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, (*mocks.LandCommodities)[0].LandArea, (*resp)[0].LandArea)
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllLandCommodity(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestLandCommodityUsecase_UpdateLandCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should update land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Update(ctx, ids.LandCommodityID, mocks.LandCommodity).Return(nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.UpdatedLandCommodity, nil).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
		assert.Equal(t, float64(200), resp.LandArea)
	})

	t.Run("should return error when validation error", func(t *testing.T) {
		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, &dto.LandCommodityUpdateDTO{LandArea: 0})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})

	t.Run("should return error when land area not enough", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(1000), nil).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Update(ctx, ids.LandCommodityID, mocks.LandCommodity).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, ids.LandCommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestLandCommodityUsecase_DeleteLandCommodity(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should delete land commodity successfully", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.LandCommodity.EXPECT().Delete(ctx, ids.LandCommodityID).Return(nil).Times(1)

		err := uc.DeleteLandCommodity(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
	})

	t.Run("should return error when land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("land commodity not found")).Times(1)

		err := uc.DeleteLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.LandCommodity.EXPECT().Delete(ctx, ids.LandCommodityID).Return(utils.NewInternalError("database error")).Times(1)

		err := uc.DeleteLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestLandCommodityUsecase_RestoreLandCommodity(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandCommodityUtils(t)

	t.Run("should restore land commodity successfully", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)
		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Restore(ctx, ids.LandCommodityID).Return(nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandCommodityID, resp.ID)
	})

	t.Run("should return error when deleted land commodity not found", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("deleted land commodity not found")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "deleted land commodity not found")
	})

	t.Run("should return error when internal error", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Restore(ctx, ids.LandCommodityID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when land area not enough", func(t *testing.T) {
		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(1000), nil).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("should return error when restored land not found", func(t *testing.T) {

		repo.LandCommodity.EXPECT().FindDeletedByID(ctx, ids.LandCommodityID).Return(mocks.LandCommodity, nil).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.LandCommodity.EXPECT().SumLandAreaByLandID(ctx, ids.LandID).Return(float64(100), nil).Times(1)

		repo.LandCommodity.EXPECT().Restore(ctx, ids.LandCommodityID).Return(nil).Times(1)

		repo.LandCommodity.EXPECT().FindByID(ctx, ids.LandCommodityID).Return(nil, utils.NewNotFoundError("restored land not found")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, ids.LandCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "restored land not found")
	})
}

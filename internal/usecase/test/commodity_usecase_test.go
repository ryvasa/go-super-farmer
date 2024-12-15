package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type CommodityRepoMock struct {
	Commodity *mock.MockCommodityRepository
}

type CommodityIDs struct {
	CommodityID uuid.UUID
}

type CommodityMocks struct {
	Commodity         *domain.Commodity
	Commodities       *[]domain.Commodity
	UpadatedCommodity *domain.Commodity
}

type CommodityDTOMock struct {
	Create *dto.CommodityCreateDTO
	Update *dto.CommodityUpdateDTO
}

func CommodityUsecaseUtils(t *testing.T) (*CommodityIDs, *CommodityMocks, *CommodityDTOMock, *CommodityRepoMock, usecase_interface.CommodityUsecase, context.Context) {
	commodityID := uuid.New()

	ids := &CommodityIDs{
		CommodityID: commodityID,
	}

	mocks := &CommodityMocks{
		Commodity: &domain.Commodity{
			ID:          commodityID,
			Name:        "test commodity",
			Description: "test commodity description",
			Code:        "12345",
		},
		Commodities: &[]domain.Commodity{
			{
				ID:          commodityID,
				Name:        "test commodity",
				Description: "test commodity description",
				Code:        "12345",
			},
		},
		UpadatedCommodity: &domain.Commodity{
			ID:          commodityID,
			Name:        "updated commodity",
			Description: "updated commodity description",
			Code:        "12345",
		},
	}

	dto := &CommodityDTOMock{
		Create: &dto.CommodityCreateDTO{
			Name:        "test commodity",
			Description: "test commodity description",
			Code:        "12345",
		},
		Update: &dto.CommodityUpdateDTO{
			Name:        "updated commodity",
			Description: "updated commodity description",
			Code:        "12345",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	uc := usecase_implementation.NewCommodityUsecase(commodityRepo)
	ctx := context.TODO()

	repo := &CommodityRepoMock{Commodity: commodityRepo}

	return ids, mocks, dto, repo, uc, ctx
}

func TestCreateCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should create commodity successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Commodity) error {
			c.ID = ids.CommodityID
			return nil
		}).Times(1)
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		resp, err := uc.CreateCommodity(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.Name, resp.Name)
		assert.Equal(t, dtos.Create.Description, resp.Description)
		assert.Equal(t, dtos.Create.Code, resp.Code)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		req := &dto.CommodityCreateDTO{}
		resp, err := uc.CreateCommodity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})
}

func TestGetAllCommodities(t *testing.T) {
	_, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)
	t.Run("should return all commodities successfully", func(t *testing.T) {

		repo.Commodity.EXPECT().FindAll(ctx).Return(mocks.Commodities, nil).Times(1)

		resp, err := uc.GetAllCommodities(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, 1)
		assert.Equal(t, (*mocks.Commodities)[0].Name, (*resp)[0].Name)
		assert.Equal(t, (*mocks.Commodities)[0].Description, (*resp)[0].Description)

	})

	t.Run("should return error when get all commodities", func(t *testing.T) {
		repo.Commodity.EXPECT().FindAll(ctx).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetAllCommodities(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetCommodityById(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)
	t.Run("should return commodity by id successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		resp, err := uc.GetCommodityById(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Commodity.Name, resp.Name)
		assert.Equal(t, mocks.Commodity.Description, resp.Description)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {

		repo.Commodity.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.GetCommodityById(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})
}

func TestUpdateCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should update commodity successfully", func(t *testing.T) {

		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Commodity.EXPECT().Update(ctx, ids.CommodityID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, c *domain.Commodity) error {
			c.ID = ids.CommodityID
			return nil
		})
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.UpadatedCommodity, nil).Times(1)

		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, dtos.Update.Name, resp.Name)
		assert.Equal(t, dtos.Update.Description, resp.Description)
		assert.Equal(t, dtos.Update.Code, resp.Code)

	})

	t.Run("should return error validation error", func(t *testing.T) {
		req := &dto.CommodityUpdateDTO{}

		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when update commodity", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Commodity.EXPECT().Update(ctx, ids.CommodityID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestDeleteCommodity(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should delete commodity successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Commodity.EXPECT().Delete(ctx, ids.CommodityID).Return(nil).Times(1)

		err := uc.DeleteCommodity(ctx, ids.CommodityID)

		assert.NoError(t, err)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		err := uc.DeleteCommodity(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when delete commodity", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Commodity.EXPECT().Delete(ctx, ids.CommodityID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteCommodity(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestRestoreCommodity(t *testing.T) {

	ids, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should restore commodity successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindDeletedByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Commodity.EXPECT().Restore(ctx, ids.CommodityID).Return(nil).Times(1)
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		resp, err := uc.RestoreCommodity(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Commodity.Name, resp.Name)
		assert.Equal(t, mocks.Commodity.Description, resp.Description)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		repo.Commodity.EXPECT().FindDeletedByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("deleted commodity not found")).Times(1)

		resp, err := uc.RestoreCommodity(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "deleted commodity not found")
	})

	t.Run("should return error when restore commodity", func(t *testing.T) {
		repo.Commodity.EXPECT().FindDeletedByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Commodity.EXPECT().Restore(ctx, ids.CommodityID).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.RestoreCommodity(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

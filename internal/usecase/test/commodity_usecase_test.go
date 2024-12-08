package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository/mock"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	uc := usecase.NewCommodityUsecase(commodityRepo)
	ctx := context.Background()

	t.Run("Test CreateCommodity successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		commodityRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Commodity) error {
			c.ID = commodityID
			return nil
		}).Times(1)
		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)

		req := &dto.CommodityCreateDTO{Name: "commodity", Description: "commodity description"}
		resp, err := uc.CreateCommodity(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
		assert.Equal(t, req.Description, resp.Description)
	})

	t.Run("Test CreateCommodity validation error", func(t *testing.T) {
		req := &dto.CommodityCreateDTO{Name: "", Description: ""}
		resp, err := uc.CreateCommodity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetAllCommodities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	uc := usecase.NewCommodityUsecase(commodityRepo)
	ctx := context.Background()

	t.Run("Test GetAllCommodities successfully", func(t *testing.T) {
		commodityID1 := uuid.New()
		commodityID2 := uuid.New()

		mockCommodity1 := &domain.Commodity{ID: commodityID1, Name: "commodity", Description: "commodity description"}
		mockCommodity2 := &domain.Commodity{ID: commodityID2, Name: "commodity", Description: "commodity description"}

		commodityRepo.EXPECT().FindAll(ctx).Return(&[]domain.Commodity{*mockCommodity1, *mockCommodity2}, nil).Times(1)

		resp, err := uc.GetAllCommodities(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, mockCommodity1.Name, (*resp)[0].Name)
		assert.Equal(t, mockCommodity1.Description, (*resp)[0].Description)
		assert.Equal(t, mockCommodity2.Name, (*resp)[1].Name)
		assert.Equal(t, mockCommodity2.Description, (*resp)[1].Description)
	})

	t.Run("Test GetAllCommodities internal error", func(t *testing.T) {
		commodityRepo.EXPECT().FindAll(ctx).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetAllCommodities(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetCommodityById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	uc := usecase.NewCommodityUsecase(commodityRepo)
	ctx := context.Background()

	t.Run("Test GetCommodityById successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)

		resp, err := uc.GetCommodityById(ctx, commodityID)

		assert.NoError(t, err)
		assert.Equal(t, mockCommodity.Name, resp.Name)
		assert.Equal(t, mockCommodity.Description, resp.Description)
	})

	t.Run("Test GetCommodityById not found", func(t *testing.T) {
		commodityID := uuid.New()

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(nil, errors.New("commodity not found")).Times(1)

		resp, err := uc.GetCommodityById(ctx, commodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test GetCommodityById internal error", func(t *testing.T) {
		commodityID := uuid.New()

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetCommodityById(ctx, commodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestUpdateCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	uc := usecase.NewCommodityUsecase(commodityRepo)
	ctx := context.Background()

	t.Run("Test UpdateCommodity successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}
		updateReq := &dto.CommodityUpdateDTO{Name: "commodity updated", Description: "commodity description updated"}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)
		commodityRepo.EXPECT().Update(ctx, commodityID, gomock.Any()).Return(nil).Times(1)
		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(&domain.Commodity{ID: commodityID, Name: "commodity updated", Description: "commodity description updated"}, nil).Times(1)

		resp, err := uc.UpdateCommodity(ctx, commodityID, updateReq)

		assert.NoError(t, err)
		assert.Equal(t, updateReq.Name, resp.Name)
		assert.Equal(t, updateReq.Description, resp.Description)
	})

	t.Run("Test UpdateCommodity validation error", func(t *testing.T) {
		commodityID := uuid.New()
		updateReq := &dto.CommodityUpdateDTO{Name: "", Description: ""}

		resp, err := uc.UpdateCommodity(ctx, commodityID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateCommodity not found", func(t *testing.T) {
		commodityID := uuid.New()
		updateReq := &dto.CommodityUpdateDTO{Name: "commodity updated", Description: "commodity description updated"}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(nil, errors.New("commodity not found")).Times(1)

		resp, err := uc.UpdateCommodity(ctx, commodityID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateCommodity internal error", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}
		updateReq := &dto.CommodityUpdateDTO{Name: "commodity updated", Description: "commodity description updated"}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)
		commodityRepo.EXPECT().Update(ctx, commodityID, gomock.Any()).Return(errors.New("internal error")).Times(1)

		resp, err := uc.UpdateCommodity(ctx, commodityID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestDeleteCommodity(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	uc := usecase.NewCommodityUsecase(commodityRepo)
	ctx := context.Background()

	t.Run("Test DeleteCommodity successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)
		commodityRepo.EXPECT().Delete(ctx, commodityID).Return(nil).Times(1)

		err := uc.DeleteCommodity(ctx, commodityID)

		assert.NoError(t, err)
	})

	t.Run("Test DeleteCommodity not found", func(t *testing.T) {
		commodityID := uuid.New()

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(nil, errors.New("commodity not found")).Times(1)

		err := uc.DeleteCommodity(ctx, commodityID)

		assert.Error(t, err)
	})

	t.Run("Test DeleteCommodity internal error", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)
		commodityRepo.EXPECT().Delete(ctx, commodityID).Return(errors.New("internal error")).Times(1)

		err := uc.DeleteCommodity(ctx, commodityID)

		assert.Error(t, err)
	})
}

func TestRestoreCommodity(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	uc := usecase.NewCommodityUsecase(commodityRepo)
	ctx := context.Background()

	t.Run("Test RestoreCommodity successfully", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		commodityRepo.EXPECT().FindDeletedByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)
		commodityRepo.EXPECT().Restore(ctx, commodityID).Return(nil).Times(1)
		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)

		resp, err := uc.RestoreCommodity(ctx, commodityID)

		assert.NoError(t, err)
		assert.Equal(t, mockCommodity.Name, resp.Name)
		assert.Equal(t, mockCommodity.Description, resp.Description)
	})

	t.Run("Test RestoreCommodity not found", func(t *testing.T) {
		commodityID := uuid.New()

		commodityRepo.EXPECT().FindDeletedByID(ctx, commodityID).Return(nil, errors.New("commodity not found")).Times(1)

		resp, err := uc.RestoreCommodity(ctx, commodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test RestoreCommodity internal error", func(t *testing.T) {
		commodityID := uuid.New()
		mockCommodity := &domain.Commodity{ID: commodityID, Name: "commodity", Description: "commodity description"}

		commodityRepo.EXPECT().FindDeletedByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)
		commodityRepo.EXPECT().Restore(ctx, commodityID).Return(errors.New("internal error")).Times(1)

		resp, err := uc.RestoreCommodity(ctx, commodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

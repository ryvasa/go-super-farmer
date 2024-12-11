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

func TestLandCommodityUsecase_CreateLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()
	t.Run("Test CreateLandCommodity successfully", func(t *testing.T) {
		mockCommodity := &domain.Commodity{ID: commodityID}
		mockLand := &domain.Land{ID: landID, LandArea: float64(200)}
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)

		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)

		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(10), nil).Times(1)

		landCommodityRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, l *domain.LandCommodity) error {
			l.ID = landCommodityID
			return nil
		}).Times(1)

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)

		req := &dto.LandCommodityCreateDTO{LandID: landID, CommodityID: commodityID, LandArea: float64(100)}
		resp, err := uc.CreateLandCommodity(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, landCommodityID, resp.ID)
	})

	t.Run("Test CreateLandCommodity land area is greater than max area", func(t *testing.T) {
		mockCommodity := &domain.Commodity{ID: commodityID}
		mockLand := &domain.Land{ID: landID, LandArea: float64(200)}

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(mockCommodity, nil).Times(1)

		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)

		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(110), nil).Times(1)

		req := &dto.LandCommodityCreateDTO{LandID: landID, CommodityID: commodityID, LandArea: float64(100)}
		resp, err := uc.CreateLandCommodity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test CreateLandCommodity validation error", func(t *testing.T) {
		commodityID := uuid.New()
		landID := uuid.New()

		resp, err := uc.CreateLandCommodity(ctx, &dto.LandCommodityCreateDTO{LandID: landID, CommodityID: commodityID, LandArea: 0})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test CreateLandCommodity notfound commodity", func(t *testing.T) {

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(nil, errors.New("commodity not found")).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, &dto.LandCommodityCreateDTO{LandID: landID, CommodityID: commodityID, LandArea: float64(100)})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("Test CreateLandCommodity notfound land", func(t *testing.T) {

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(&domain.Commodity{ID: commodityID}, nil).Times(1)

		landRepo.EXPECT().FindByID(ctx, landID).Return(nil, errors.New("land not found")).Times(1)

		resp, err := uc.CreateLandCommodity(ctx, &dto.LandCommodityCreateDTO{LandID: landID, CommodityID: commodityID, LandArea: float64(100)})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})
}

func TestLandCommodityUsecase_GetLandCommodityByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetLandCommodityByID not found", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(nil, errors.New("land commodity not found")).Times(1)

		resp, err := uc.GetLandCommodityByID(ctx, landCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land commodity not found")
	})
	t.Run("Test GetLandCommodity successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)

		resp, err := uc.GetLandCommodityByID(ctx, landCommodityID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, landCommodityID, resp.ID)
		assert.Equal(t, commodityID, resp.CommodityID)
		assert.Equal(t, landID, resp.LandID)
		assert.Equal(t, float64(100), resp.LandArea)
	})
}

func TestLandCommodityUsecase_GetLandCommodityByLandID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetLandCommodityByLandID successfully", func(t *testing.T) {
		mockLandCommodity1 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLandCommodity2 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}

		landCommodityRepo.EXPECT().FindByLandID(ctx, landID).Return(&[]domain.LandCommodity{*mockLandCommodity1, *mockLandCommodity2}, nil).Times(1)

		resp, err := uc.GetLandCommodityByLandID(ctx, landID)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, mockLandCommodity1.LandArea, (*resp)[0].LandArea)
		assert.Equal(t, mockLandCommodity1.CommodityID, (*resp)[0].CommodityID)
		assert.Equal(t, mockLandCommodity2.LandArea, (*resp)[1].LandArea)
		assert.Equal(t, mockLandCommodity2.CommodityID, (*resp)[1].CommodityID)
	})

	t.Run("Test GetLandCommodityByLandID internal error", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindByLandID(ctx, landID).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetLandCommodityByLandID(ctx, landID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestLandCommodityUsecase_GetLandCommodityByCommodityID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetLandCommodityByCommodityID successfully", func(t *testing.T) {
		mockLandCommodity1 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLandCommodity2 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}

		landCommodityRepo.EXPECT().FindByCommodityID(ctx, commodityID).Return(&[]domain.LandCommodity{*mockLandCommodity1, *mockLandCommodity2}, nil).Times(1)

		resp, err := uc.GetLandCommodityByCommodityID(ctx, commodityID)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, mockLandCommodity1.LandArea, (*resp)[0].LandArea)
		assert.Equal(t, mockLandCommodity1.CommodityID, (*resp)[0].CommodityID)
		assert.Equal(t, mockLandCommodity2.LandArea, (*resp)[1].LandArea)
		assert.Equal(t, mockLandCommodity2.CommodityID, (*resp)[1].CommodityID)
	})

	t.Run("Test GetLandCommodityByCommodityID internal error", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindByCommodityID(ctx, commodityID).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetLandCommodityByCommodityID(ctx, commodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestLandCommodityUsecase_GetAllLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test GetAllLandCommodity successfully", func(t *testing.T) {
		mockLandCommodity1 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLandCommodity2 := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}

		landCommodityRepo.EXPECT().FindAll(ctx).Return(&[]domain.LandCommodity{*mockLandCommodity1, *mockLandCommodity2}, nil).Times(1)

		resp, err := uc.GetAllLandCommodity(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, mockLandCommodity1.LandArea, (*resp)[0].LandArea)
		assert.Equal(t, mockLandCommodity1.CommodityID, (*resp)[0].CommodityID)
		assert.Equal(t, mockLandCommodity2.LandArea, (*resp)[1].LandArea)
		assert.Equal(t, mockLandCommodity2.CommodityID, (*resp)[1].CommodityID)
	})

	t.Run("Test GetAllLandCommodity internal error", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindAll(ctx).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetAllLandCommodity(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestLandCommodityUsecase_UpdateLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test UpdateLandCommodity successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLand := &domain.Land{ID: landID, LandArea: float64(1000)}

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)
		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(&domain.Commodity{ID: commodityID}, nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)
		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(100), nil).Times(1)

		landCommodityRepo.EXPECT().Update(ctx, landCommodityID, mockLandCommodity).Return(nil).Times(1)

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(&domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}, nil).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx,
			landCommodityID, &dto.LandCommodityUpdateDTO{LandArea: float64(200), CommodityID: commodityID, LandID: landID})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, landCommodityID, resp.ID)
		assert.Equal(t, float64(200), resp.LandArea)
	})

	t.Run("Test UpdateLandCommodity validation error", func(t *testing.T) {

		resp, err := uc.UpdateLandCommodity(ctx, landCommodityID, &dto.LandCommodityUpdateDTO{LandArea: 0})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateLandCommodity notfound", func(t *testing.T) {

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(nil, errors.New("land commodity not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, landCommodityID, &dto.LandCommodityUpdateDTO{LandArea: float64(200)})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("Test UpdateLandCommodity notfound commodity", func(t *testing.T) {

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(&domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}, nil).Times(1)

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(nil, errors.New("commodity not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, landCommodityID, &dto.LandCommodityUpdateDTO{LandArea: float64(200), CommodityID: commodityID, LandID: landID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("Test UpdateLandCommodity notfound land", func(t *testing.T) {

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(&domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}, nil).Times(1)

		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(&domain.Commodity{ID: commodityID}, nil).Times(1)

		landRepo.EXPECT().FindByID(ctx, landID).Return(nil, errors.New("land not found")).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx, landCommodityID, &dto.LandCommodityUpdateDTO{LandArea: float64(200), CommodityID: commodityID, LandID: landID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})

	t.Run("Test UpdateLandCommodity land area not enough", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLand := &domain.Land{ID: landID, LandArea: float64(100)}

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)
		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(&domain.Commodity{ID: commodityID}, nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)
		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(100), nil).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx,
			landCommodityID, &dto.LandCommodityUpdateDTO{LandArea: float64(200), CommodityID: commodityID, LandID: landID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("Test UpdateLandCommodity internal error", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}
		mockLand := &domain.Land{ID: landID, LandArea: float64(1000)}

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)
		commodityRepo.EXPECT().FindByID(ctx, commodityID).Return(&domain.Commodity{ID: commodityID}, nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)
		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(100), nil).Times(1)

		landCommodityRepo.EXPECT().Update(ctx, landCommodityID, mockLandCommodity).Return(errors.New("internal error")).Times(1)

		// landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(&domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(200)}, nil).Times(1)

		resp, err := uc.UpdateLandCommodity(ctx,
			landCommodityID, &dto.LandCommodityUpdateDTO{LandArea: float64(200), CommodityID: commodityID, LandID: landID})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestLandCommodityUsecase_DeleteLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test DeleteLandCommodity successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)

		landCommodityRepo.EXPECT().Delete(ctx, landCommodityID).Return(nil).Times(1)

		err := uc.DeleteLandCommodity(ctx, landCommodityID)

		assert.NoError(t, err)
	})

	t.Run("Test DeleteLandCommodity notfound", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(nil, errors.New("land commodity not found")).Times(1)

		err := uc.DeleteLandCommodity(ctx, landCommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "land commodity not found")
	})

	t.Run("Test DeleteLandCommodity internal error", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(&domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}, nil).Times(1)

		landCommodityRepo.EXPECT().Delete(ctx, landCommodityID).Return(errors.New("database error")).Times(1)

		err := uc.DeleteLandCommodity(ctx, landCommodityID)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestLandCommodityUsecase_RestoreLandCommodity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock.NewMockCommodityRepository(ctrl)
	landCommodityRepo := mock.NewMockLandCommodityRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandCommodityUsecase(landCommodityRepo, landRepo, commodityRepo)
	ctx := context.Background()

	landCommodityID := uuid.New()
	commodityID := uuid.New()
	landID := uuid.New()

	t.Run("Test RestoreLandCommodity successfully", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		landCommodityRepo.EXPECT().FindDeletedByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: float64(1000)}, nil).Times(1)
		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(100), nil).Times(1)

		landCommodityRepo.EXPECT().Restore(ctx, landCommodityID).Return(nil).Times(1)

		landCommodityRepo.EXPECT().FindByID(ctx, landCommodityID).Return(&domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}, nil).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, landCommodityID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, landCommodityID, resp.ID)
	})

	t.Run("Test RestoreLandCommodity notfound", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindDeletedByID(ctx, landCommodityID).Return(nil, errors.New("land commodity not found")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, landCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "deleted land commodity not found")
	})

	t.Run("Test RestoreLandCommodity internal error", func(t *testing.T) {
		landCommodityRepo.EXPECT().FindDeletedByID(ctx, landCommodityID).Return(&domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}, nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: float64(1000)}, nil).Times(1)
		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(100), nil).Times(1)
		landCommodityRepo.EXPECT().Restore(ctx, landCommodityID).Return(errors.New("database error")).Times(1)

		resp, err := uc.RestoreLandCommodity(ctx, landCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "database error")
	})

	t.Run("Test RestoreLandCommodity land area not enough", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		landCommodityRepo.EXPECT().FindDeletedByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: float64(100)}, nil).Times(1)
		landCommodityRepo.EXPECT().SumLandAreaByLandID(ctx, landID).Return(float64(100), nil).Times(1)
		resp, err := uc.RestoreLandCommodity(ctx, landCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land area not enough")
	})

	t.Run("Test RestoreLandCommodity land not found", func(t *testing.T) {
		mockLandCommodity := &domain.LandCommodity{ID: landCommodityID, CommodityID: commodityID, LandID: landID, LandArea: float64(100)}

		landCommodityRepo.EXPECT().FindDeletedByID(ctx, landCommodityID).Return(mockLandCommodity, nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(nil, errors.New("land not found")).Times(1)
		resp, err := uc.RestoreLandCommodity(ctx, landCommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})
}

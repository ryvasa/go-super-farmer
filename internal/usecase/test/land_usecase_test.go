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
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandUsecase(landRepo, userRepo)
	ctx := context.Background()

	t.Run("Test CreateLand, successfully", func(t *testing.T) {
		landID := uuid.New()
		userID := uuid.New()
		mockLand := &domain.Land{ID: landID, LandArea: 100, Certificate: "cert", UserID: userID}

		landRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, l *domain.Land) error {
			l.ID = landID
			return nil
		}).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)

		req := &dto.LandCreateDTO{LandArea: 100, Certificate: "cert"}
		resp, err := uc.CreateLand(ctx, userID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.LandArea, resp.LandArea)
		assert.Equal(t, req.Certificate, resp.Certificate)
	})

	t.Run("Test CreateLand, validation error", func(t *testing.T) {
		userID := uuid.New()

		req := &dto.LandCreateDTO{LandArea: 0, Certificate: ""}
		resp, err := uc.CreateLand(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetLandByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandUsecase(landRepo, nil)
	ctx := context.Background()

	t.Run("Test GetLandByID successfully", func(t *testing.T) {
		landID := uuid.New()
		userID := uuid.New()
		mockLand := &domain.Land{ID: landID, LandArea: 100, Certificate: "cert", UserID: userID}

		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)

		resp, err := uc.GetLandByID(ctx, landID)

		assert.NoError(t, err)
		assert.Equal(t, mockLand.LandArea, resp.LandArea)
		assert.Equal(t, mockLand.Certificate, resp.Certificate)
	})

	t.Run("Test GetLandByID not found", func(t *testing.T) {
		landID := uuid.New()

		landRepo.EXPECT().FindByID(ctx, landID).Return(nil, errors.New("land not found")).Times(1)

		resp, err := uc.GetLandByID(ctx, landID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetLandByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landRepo := mock.NewMockLandRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	uc := usecase.NewLandUsecase(landRepo, userRepo)
	ctx := context.Background()

	t.Run("Test GetLandByUserID successfully", func(t *testing.T) {
		userID := uuid.New()
		landID1 := uuid.New()
		landID2 := uuid.New()

		mockLand1 := &domain.Land{ID: landID1, LandArea: 100, Certificate: "cert", UserID: userID}
		mockLand2 := &domain.Land{ID: landID2, LandArea: 200, Certificate: "cert2", UserID: userID}

		userRepo.EXPECT().FindByID(ctx, userID).Return(&domain.User{ID: userID}, nil).Times(1)
		landRepo.EXPECT().FindByUserID(ctx, userID).Return(&[]domain.Land{*mockLand1, *mockLand2}, nil).Times(1)

		resp, err := uc.GetLandByUserID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, *resp, 2)
		assert.Equal(t, mockLand1.LandArea, (*resp)[0].LandArea)
		assert.Equal(t, mockLand1.Certificate, (*resp)[0].Certificate)
		assert.Equal(t, mockLand2.LandArea, (*resp)[1].LandArea)
		assert.Equal(t, mockLand2.Certificate, (*resp)[1].Certificate)
	})

	t.Run("Test GetLandByUserID user not found", func(t *testing.T) {
		userID := uuid.New()

		userRepo.EXPECT().FindByID(ctx, userID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		resp, err := uc.GetLandByUserID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("Test GetLandByUserID internal error on FindByUserID", func(t *testing.T) {
		userID := uuid.New()

		userRepo.EXPECT().FindByID(ctx, userID).Return(&domain.User{ID: userID}, nil).Times(1)
		landRepo.EXPECT().FindByUserID(ctx, userID).Return(nil, errors.New("database error")).Times(1)

		resp, err := uc.GetLandByUserID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "database error")
	})
}

func TestGetAllLands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandUsecase(landRepo, nil)
	ctx := context.Background()

	t.Run("Test GetAllLands successfully", func(t *testing.T) {
		landID1 := uuid.New()
		landID2 := uuid.New()

		mockLand1 := &domain.Land{ID: landID1, LandArea: 100, Certificate: "cert"}
		mockLand2 := &domain.Land{ID: landID2, LandArea: 100, Certificate: "cert"}

		landRepo.EXPECT().FindAll(ctx).Return(&[]domain.Land{*mockLand1, *mockLand2}, nil).Times(1)

		resp, err := uc.GetAllLands(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, 2)
		assert.Equal(t, mockLand1.LandArea, (*resp)[0].LandArea)
		assert.Equal(t, mockLand1.Certificate, (*resp)[0].Certificate)
		assert.Equal(t, mockLand2.LandArea, (*resp)[1].LandArea)
		assert.Equal(t, mockLand2.Certificate, (*resp)[1].Certificate)
	})

	t.Run("Test GetAllLands internal error", func(t *testing.T) {
		landRepo.EXPECT().FindAll(ctx).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetAllLands(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestUpdateLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandUsecase(landRepo, nil)
	ctx := context.Background()

	t.Run("Test UpdateLand successfully", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		mockLand := &domain.Land{ID: landID, LandArea: 100, Certificate: "old_cert", UserID: userID}
		updateReq := &dto.LandUpdateDTO{LandArea: 200, Certificate: "new_cert"}

		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)
		landRepo.EXPECT().Update(ctx, landID, mockLand).Return(nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: 200, Certificate: "new_cert", UserID: userID}, nil).Times(1)

		resp, err := uc.UpdateLand(ctx, userID, landID, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, updateReq.LandArea, resp.LandArea)
		assert.Equal(t, updateReq.Certificate, resp.Certificate)
	})

	t.Run("Test UpdateLand not found", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		updateReq := &dto.LandUpdateDTO{LandArea: 200, Certificate: "new_cert"}

		landRepo.EXPECT().FindByID(ctx, landID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.UpdateLand(ctx, userID, landID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})

	t.Run("Test UpdateLand validation error", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		updateReq := &dto.LandUpdateDTO{LandArea: -1, Certificate: ""}

		resp, err := uc.UpdateLand(ctx, userID, landID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Test UpdateLand internal error on update", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		mockLand := &domain.Land{ID: landID, LandArea: 100, Certificate: "old_cert", UserID: userID}
		updateReq := &dto.LandUpdateDTO{LandArea: 200, Certificate: "new_cert"}

		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)
		landRepo.EXPECT().Update(ctx, landID, mockLand).Return(errors.New("database error")).Times(1)

		resp, err := uc.UpdateLand(ctx, userID, landID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "database error")
	})

	t.Run("Test UpdateLand internal error on FindByID after update", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()
		mockLand := &domain.Land{ID: landID, LandArea: 100, Certificate: "old_cert", UserID: userID}
		updateReq := &dto.LandUpdateDTO{LandArea: 200, Certificate: "new_cert"}

		landRepo.EXPECT().FindByID(ctx, landID).Return(mockLand, nil).Times(1)
		landRepo.EXPECT().Update(ctx, landID, mockLand).Return(nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(nil, errors.New("database error")).Times(1)

		resp, err := uc.UpdateLand(ctx, userID, landID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "database error")
	})
}

func TestDeleteLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandUsecase(landRepo, nil)
	ctx := context.Background()

	t.Run("Test DeleteLand successfully", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()

		landRepo.EXPECT().FindByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: 100, Certificate: "cert", UserID: userID}, nil).Times(1)
		landRepo.EXPECT().Delete(ctx, landID).Return(nil).Times(1)

		err := uc.DeleteLand(ctx, landID)

		assert.NoError(t, err)
	})

	t.Run("Test DeleteLand not found", func(t *testing.T) {
		landID := uuid.New()

		landRepo.EXPECT().FindByID(ctx, landID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		err := uc.DeleteLand(ctx, landID)

		assert.Error(t, err)
		assert.EqualError(t, err, "land not found")
	})

	t.Run("Test DeleteLand internal error", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()

		landRepo.EXPECT().FindByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: 100, Certificate: "cert", UserID: userID}, nil).Times(1)
		landRepo.EXPECT().Delete(ctx, landID).Return(errors.New("database error")).Times(1)

		err := uc.DeleteLand(ctx, landID)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestRestoreLand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landRepo := mock.NewMockLandRepository(ctrl)
	uc := usecase.NewLandUsecase(landRepo, nil)
	ctx := context.Background()

	t.Run("Test RestoreLand successfully", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()

		landRepo.EXPECT().FindDeletedByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: 100, Certificate: "cert", UserID: userID}, nil).Times(1)
		landRepo.EXPECT().Restore(ctx, landID).Return(nil).Times(1)
		landRepo.EXPECT().FindByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: 100, Certificate: "cert", UserID: userID}, nil).Times(1)

		resp, err := uc.RestoreLand(ctx, landID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, landID, resp.ID)
	})

	t.Run("Test RestoreLand not found", func(t *testing.T) {
		landID := uuid.New()

		landRepo.EXPECT().FindDeletedByID(ctx, landID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.RestoreLand(ctx, landID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "deleted land not found")
	})

	t.Run("Test RestoreLand internal error", func(t *testing.T) {
		userID := uuid.New()
		landID := uuid.New()

		landRepo.EXPECT().FindDeletedByID(ctx, landID).Return(&domain.Land{ID: landID, LandArea: 100, Certificate: "cert", UserID: userID}, nil).Times(1)
		landRepo.EXPECT().Restore(ctx, landID).Return(errors.New("database error")).Times(1)

		resp, err := uc.RestoreLand(ctx, landID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "database error")
	})
}

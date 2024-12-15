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
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type LandRepoMock struct {
	Land *mock.MockLandRepository
	User *mock.MockUserRepository
}

type LandIDs struct {
	LandID uuid.UUID
	UserID uuid.UUID
}

type LandMocks struct {
	Land        *domain.Land
	Lands       *[]domain.Land
	UpdatedLand *domain.Land
	User        *domain.User
}

type LandDTOMock struct {
	Create *dto.LandCreateDTO
	Update *dto.LandUpdateDTO
}

func LandUsecaseUtils(t *testing.T) (*LandIDs, *LandMocks, *LandDTOMock, *LandRepoMock, usecase_interface.LandUsecase, context.Context) {
	landID := uuid.New()
	userID := uuid.New()

	ids := &LandIDs{
		LandID: landID,
		UserID: userID,
	}

	mocks := &LandMocks{
		Land: &domain.Land{
			ID:          landID,
			LandArea:    100,
			Certificate: "cert",
			UserID:      userID,
		},
		Lands: &[]domain.Land{
			{
				ID:          landID,
				LandArea:    100,
				Certificate: "cert",
				UserID:      userID,
			},
		},
		UpdatedLand: &domain.Land{
			ID:          landID,
			LandArea:    99,
			Certificate: "updated cert",
			UserID:      userID,
		},
		User: &domain.User{
			ID: userID,
		},
	}

	dto := &LandDTOMock{
		Create: &dto.LandCreateDTO{
			LandArea:    100,
			Certificate: "cert",
		},
		Update: &dto.LandUpdateDTO{
			LandArea:    99,
			Certificate: "updated cert",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	landRepo := mock.NewMockLandRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	uc := usecase_implementation.NewLandUsecase(landRepo, userRepo)
	ctx := context.TODO()

	repo := &LandRepoMock{Land: landRepo, User: userRepo}

	return ids, mocks, dto, repo, uc, ctx
}

func TestCreateLand(t *testing.T) {

	ids, mocks, dtos, repo, uc, ctx := LandUsecaseUtils(t)

	t.Run("should create land successfully", func(t *testing.T) {

		repo.Land.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, l *domain.Land) error {
			l.ID = ids.LandID
			return nil
		}).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		resp, err := uc.CreateLand(ctx, ids.UserID, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.LandArea, resp.LandArea)
		assert.Equal(t, dtos.Create.Certificate, resp.Certificate)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		userID := uuid.New()

		req := &dto.LandCreateDTO{LandArea: 0, Certificate: ""}
		resp, err := uc.CreateLand(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error internal error", func(t *testing.T) {
		userID := uuid.New()
		repo.Land.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateLand(ctx, userID, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when find created land ", func(t *testing.T) {
		repo.Land.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, l *domain.Land) error {
			l.ID = ids.LandID
			return nil
		}).Times(1)
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateLand(ctx, ids.UserID, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetLandByID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandUsecaseUtils(t)

	t.Run("should return land successfully", func(t *testing.T) {

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		resp, err := uc.GetLandByID(ctx, ids.LandID)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Land.LandArea, resp.LandArea)
		assert.Equal(t, mocks.Land.Certificate, resp.Certificate)
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.GetLandByID(ctx, ids.LandID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})
}

func TestGetLandByUserID(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandUsecaseUtils(t)

	t.Run("should return lands successfully", func(t *testing.T) {

		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)

		repo.Land.EXPECT().FindByUserID(ctx, ids.UserID).Return(mocks.Lands, nil).Times(1)

		resp, err := uc.GetLandByUserID(ctx, ids.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, *resp, 1)
		assert.Equal(t, (*mocks.Lands)[0].LandArea, (*resp)[0].LandArea)
		assert.Equal(t, (*mocks.Lands)[0].Certificate, (*resp)[0].Certificate)
	})

	t.Run("should return error user not found", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(nil, utils.NewNotFoundError("user not found")).Times(1)

		resp, err := uc.GetLandByUserID(ctx, ids.UserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("should return error internal error when find lands", func(t *testing.T) {
		repo.User.EXPECT().FindByID(ctx, ids.UserID).Return(mocks.User, nil).Times(1)
		repo.Land.EXPECT().FindByUserID(ctx, ids.UserID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.GetLandByUserID(ctx, ids.UserID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestGetAllLands(t *testing.T) {
	_, mocks, _, repo, uc, ctx := LandUsecaseUtils(t)

	t.Run("should return lands successfully", func(t *testing.T) {
		repo.Land.EXPECT().FindAll(ctx).Return(mocks.Lands, nil).Times(1)
		resp, err := uc.GetAllLands(ctx)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Lands, resp)
		assert.Len(t, (*mocks.Lands), 1)
		assert.Equal(t, (*mocks.Lands)[0].LandArea, (*resp)[0].LandArea)
		assert.Equal(t, (*mocks.Lands)[0].Certificate, (*resp)[0].Certificate)
	})

	t.Run("should return error internal error", func(t *testing.T) {
		repo.Land.EXPECT().FindAll(ctx).Return(nil, errors.New("internal error")).Times(1)

		resp, err := uc.GetAllLands(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

func TestUpdateLand(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := LandUsecaseUtils(t)
	t.Run("should update land successfully", func(t *testing.T) {

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.Land.EXPECT().Update(ctx, ids.LandID, mocks.UpdatedLand).DoAndReturn(func(ctx context.Context, id uuid.UUID, l *domain.Land) error {
			l.ID = ids.LandID
			return nil
		}).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.UpdatedLand, nil).Times(1)

		resp, err := uc.UpdateLand(ctx, ids.UserID, ids.LandID, dtos.Update)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, dtos.Update.LandArea, resp.LandArea)
		assert.Equal(t, dtos.Update.Certificate, resp.Certificate)
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.UpdateLand(ctx, ids.UserID, ids.LandID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "land not found")
	})

	t.Run("should return error validation error", func(t *testing.T) {
		updateReq := &dto.LandUpdateDTO{LandArea: -1, Certificate: ""}

		resp, err := uc.UpdateLand(ctx, ids.UserID, ids.LandID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error internal error on update", func(t *testing.T) {
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.Land.EXPECT().Update(ctx, ids.LandID, mocks.Land).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.UpdateLand(ctx, ids.UserID, ids.LandID, dtos.Update)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error internal error when find updated land", func(t *testing.T) {
		// Setup ekspektasi
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.Land.EXPECT().Update(ctx, ids.LandID, mocks.Land).DoAndReturn(func(ctx context.Context, id uuid.UUID, l *domain.Land) error {
			l.ID = ids.LandID
			return nil
		}).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewInternalError("internal error")).Times(1)

		// Eksekusi fungsi
		resp, err := uc.UpdateLand(ctx, ids.UserID, ids.LandID, dtos.Update)

		// Validasi hasil
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

}

func TestDeleteLand(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandUsecaseUtils(t)

	t.Run("should delete land successfully", func(t *testing.T) {

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)
		repo.Land.EXPECT().Delete(ctx, ids.LandID).Return(nil).Times(1)

		err := uc.DeleteLand(ctx, ids.LandID)

		assert.NoError(t, err)
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)
		err := uc.DeleteLand(ctx, ids.LandID)

		assert.Error(t, err)
		assert.EqualError(t, err, "land not found")
	})

	t.Run("should return error internal error", func(t *testing.T) {
		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)
		repo.Land.EXPECT().Delete(ctx, ids.LandID).Return(utils.NewInternalError("internal error")).Times(1)

		err := uc.DeleteLand(ctx, ids.LandID)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})
}

func TestRestoreLand(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := LandUsecaseUtils(t)

	t.Run("should restore land successfully", func(t *testing.T) {

		repo.Land.EXPECT().FindDeletedByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.Land.EXPECT().Restore(ctx, ids.LandID).Return(nil).Times(1)

		repo.Land.EXPECT().FindByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		resp, err := uc.RestoreLand(ctx, ids.LandID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, ids.LandID, resp.ID)
	})

	t.Run("should return error when land not found", func(t *testing.T) {
		repo.Land.EXPECT().FindDeletedByID(ctx, ids.LandID).Return(nil, utils.NewNotFoundError("land not found")).Times(1)

		resp, err := uc.RestoreLand(ctx, ids.LandID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "deleted land not found")
	})

	t.Run("Test RestoreLand internal error", func(t *testing.T) {

		repo.Land.EXPECT().FindDeletedByID(ctx, ids.LandID).Return(mocks.Land, nil).Times(1)

		repo.Land.EXPECT().Restore(ctx, ids.LandID).Return(errors.New("internal error")).Times(1)

		resp, err := uc.RestoreLand(ctx, ids.LandID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})
}

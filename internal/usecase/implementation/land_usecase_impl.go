package usecase_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type LandUsecaseImpl struct {
	landRepo repository_interface.LandRepository
	userRepo repository_interface.UserRepository
}

func NewLandUsecase(landRepo repository_interface.LandRepository, userRepo repository_interface.UserRepository) usecase_interface.LandUsecase {
	return &LandUsecaseImpl{landRepo, userRepo}
}

func (u *LandUsecaseImpl) CreateLand(ctx context.Context, userId uuid.UUID, req *dto.LandCreateDTO) (*domain.Land, error) {
	land := domain.Land{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	land.LandArea = req.LandArea
	land.Certificate = req.Certificate
	land.CityID = req.CityID
	land.UserID = userId
	land.ID = uuid.New()

	err := u.landRepo.Create(ctx, &land)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	createdLand, err := u.landRepo.FindByID(ctx, land.ID)

	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdLand, nil
}

func (u *LandUsecaseImpl) GetLandByID(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	land, err := u.landRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("land not found")
	}
	return land, nil
}

func (u *LandUsecaseImpl) GetLandByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Land, error) {
	_, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, utils.NewNotFoundError("user not found")
	}
	land, err := u.landRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return land, nil
}

func (u *LandUsecaseImpl) GetAllLands(ctx context.Context) ([]*domain.Land, error) {
	lands, err := u.landRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return lands, nil
}

func (u *LandUsecaseImpl) UpdateLand(ctx context.Context, userId, id uuid.UUID, req *dto.LandUpdateDTO) (*domain.Land, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	land, err := u.landRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("land not found")
	}
	land.LandArea = req.LandArea
	land.Certificate = req.Certificate

	err = u.landRepo.Update(ctx, id, land)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedLand, err := u.landRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedLand, nil
}

func (u *LandUsecaseImpl) DeleteLand(ctx context.Context, id uuid.UUID) error {
	_, err := u.landRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("land not found")
	}

	err = u.landRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	return nil
}

func (u *LandUsecaseImpl) RestoreLand(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	_, err := u.landRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted land not found")
	}

	err = u.landRepo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	restoredLand, err := u.landRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return restoredLand, nil
}

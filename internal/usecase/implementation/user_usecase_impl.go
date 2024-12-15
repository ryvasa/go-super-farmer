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

type UserUsecaseImpl struct {
	repo repository_interface.UserRepository
	hash utils.Hasher
}

func NewUserUsecase(repo repository_interface.UserRepository, hash utils.Hasher) usecase_interface.UserUsecase {
	return &UserUsecaseImpl{repo, hash}
}

func (uc *UserUsecaseImpl) Register(ctx context.Context, req *dto.UserCreateDTO) (*dto.UserResponseDTO, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	user.Name = req.Name
	user.Email = req.Email
	user.ID = uuid.New()

	hashedPassword, err := uc.hash.HashPassword(req.Password)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	user.Password = hashedPassword
	err = uc.repo.Create(ctx, &user)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	createdUser, err := uc.repo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return utils.UserDtoFormat(createdUser), nil
}

func (uc *UserUsecaseImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return utils.UserDtoFormat(user), nil
}

func (uc *UserUsecaseImpl) GetAllUsers(ctx context.Context) (*[]dto.UserResponseDTO, error) {
	users, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	usersDto := make([]dto.UserResponseDTO, 0)
	for _, user := range *users {
		usersDto = append(usersDto, *utils.UserDtoFormat(&user))
	}
	return &usersDto, nil
}

func (uc *UserUsecaseImpl) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Phone = &req.Phone

	if req.Password != "" {
		hashedPassword, err := uc.hash.HashPassword(req.Password)
		if err != nil {
			return nil, utils.NewInternalError(err.Error())
		}
		user.Password = hashedPassword
	}

	err = uc.repo.Update(ctx, id, &user)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	updatedUser, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return utils.UserDtoFormat(updatedUser), nil
}

func (uc *UserUsecaseImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError(err.Error())
	}
	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (uc *UserUsecaseImpl) RestoreUser(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error) {
	_, err := uc.repo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	err = uc.repo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	restoredUser, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return utils.UserDtoFormat(restoredUser), err
}
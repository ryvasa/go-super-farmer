package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type UserUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{repo: repo}
}

func (uc *UserUsecaseImpl) Register(req *dto.UserCreateDTO) (*dto.UserResponseDTO, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return &dto.UserResponseDTO{}, utils.NewValidationError(err)
	}
	user.Name = req.Name
	user.Email = req.Email

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return &dto.UserResponseDTO{}, utils.NewInternalError(err.Error())
	}

	user.Password = hashedPassword
	err = uc.repo.Create(&user)
	if err != nil {
		return &dto.UserResponseDTO{}, err
	}
	createdUser, err := uc.repo.FindById(user.ID)
	if err != nil {
		return &dto.UserResponseDTO{}, err
	}

	return utils.UserDtoFormat(createdUser), nil
}

func (uc *UserUsecaseImpl) GetUserByID(id int64) (*dto.UserResponseDTO, error) {
	user, err := uc.repo.FindById(id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return utils.UserDtoFormat(user), nil
}

func (uc *UserUsecaseImpl) GetAllUsers() (*[]dto.UserResponseDTO, error) {
	users, err := uc.repo.FindAll()
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	usersDto := make([]dto.UserResponseDTO, 0)
	for _, user := range *users {
		usersDto = append(usersDto, *utils.UserDtoFormat(&user))
	}
	return &usersDto, nil
}

func (uc *UserUsecaseImpl) UpdateUser(id int64, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Password = req.Password
	err := uc.repo.Update(id, &user)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	updatedUser, err := uc.repo.FindById(id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return utils.UserDtoFormat(updatedUser), nil
}

func (uc *UserUsecaseImpl) DeleteUser(id int64) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (uc *UserUsecaseImpl) RestoreUser(id int64) (*dto.UserResponseDTO, error) {
	_, err := uc.repo.FindDeletedById(id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	err = uc.repo.Restore(id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	restoredUser, err := uc.repo.FindById(id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return utils.UserDtoFormat(restoredUser), err
}

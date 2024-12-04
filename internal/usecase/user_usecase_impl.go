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

func (uc *UserUsecaseImpl) Register(req *dto.UserCreateDTO) (*domain.User, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return &user, utils.NewValidationError(err)
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Password = req.Password
	err := uc.repo.Create(&user)
	if err != nil {
		return &user, err
	}
	createdUser, err := uc.repo.FindById(user.ID)
	if err != nil {
		return &user, err
	}
	return createdUser, nil
}

func (uc *UserUsecaseImpl) GetUserByID(id int64) (*domain.User, error) {
	user, err := uc.repo.FindById(id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return user, nil
}

func (uc *UserUsecaseImpl) GetAllUsers() (*[]domain.User, error) {
	users, err := uc.repo.FindAll()
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return users, nil
}

func (uc *UserUsecaseImpl) UpdateUser(id int64, req *dto.UserUpdateDTO) (*domain.User, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return &user, utils.NewValidationError(err)
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Password = req.Password
	err := uc.repo.Update(id, &user)
	if err != nil {
		return &user, utils.NewInternalError(err.Error())
	}
	updatedUser, err := uc.repo.FindById(id)
	if err != nil {
		return &user, utils.NewNotFoundError(err.Error())
	}
	return updatedUser, nil
}

func (uc *UserUsecaseImpl) DeleteUser(id int64) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (uc *UserUsecaseImpl) RestoreUser(id int64) (*domain.User, error) {
	user := domain.User{}
	_, err := uc.repo.FindDeletedById(id)
	if err != nil {
		return &user, utils.NewNotFoundError(err.Error())
	}
	err = uc.repo.Restore(id)
	if err != nil {
		return &user, utils.NewInternalError(err.Error())
	}
	restoredUser, err := uc.repo.FindById(id)
	if err != nil {
		return &user, err
	}
	return restoredUser, err
}

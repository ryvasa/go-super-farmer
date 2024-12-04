package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
)

type UserUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{repo: repo}
}

func (uc *UserUsecaseImpl) Register(user *domain.User) error {
	return uc.repo.Create(user)
}

func (uc *UserUsecaseImpl) GetUserByID(id int64) (*domain.User, error) {
	return uc.repo.FindById(id)
}

func (uc *UserUsecaseImpl) GetAllUsers() ([]domain.User, error) {
	return uc.repo.FindAll()
}

func (uc *UserUsecaseImpl) UpdateUser(id int64, user *domain.User) error {
	return uc.repo.Update(id, user)
}

func (uc *UserUsecaseImpl) DeleteUser(id int64) error {
	return uc.repo.Delete(id)
}

func (uc *UserUsecaseImpl) RestoreUser(id int64) error {
	return uc.repo.Restore(id)
}

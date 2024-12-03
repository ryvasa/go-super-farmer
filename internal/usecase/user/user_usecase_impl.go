package usecase

import (
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository "github.com/ryvasa/go-super-farmer/internal/repository/user"
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

func (uc *UserUsecaseImpl) GetUserByID(id uint) (*domain.User, error) {
	return uc.repo.FindById(id)
}

package usecase

import "github.com/ryvasa/go-super-farmer/internal/domain"

type UserUsecase struct {
	repo UserRepository
}

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) Register(user *domain.User) error {
	return uc.repo.Create(user)
}

func (uc *UserUsecase) GetUserByID(id uint) (*domain.User, error) {
	return uc.repo.FindByID(id)
}

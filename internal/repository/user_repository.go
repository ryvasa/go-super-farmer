package repository

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint64) (*domain.User, error)
	FindAll(ctx context.Context) (*[]domain.User, error)
	Delete(ctx context.Context, id uint64) error
	Restore(ctx context.Context, id uint64) error
	Update(ctx context.Context, id uint64, user *domain.User) error
	FindDeletedByID(ctx context.Context, id uint64) (*domain.User, error)
}

package repository_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type ProvinceRepository interface {
	Create(ctx context.Context, province *domain.Province) error
	FindByID(ctx context.Context, id int64) (*domain.Province, error)
	FindAll(ctx context.Context) ([]*domain.Province, error)
	Update(ctx context.Context, id int64, province *domain.Province) error
	Delete(ctx context.Context, id int64) error
}

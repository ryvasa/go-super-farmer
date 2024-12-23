package repository_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
)

type CityRepository interface {
	Create(ctx context.Context, city *domain.City) error
	FindByID(ctx context.Context, id int64) (*domain.City, error)
	FindAll(ctx context.Context) ([]*domain.City, error)
	Update(ctx context.Context, id int64, city *domain.City) error
	Delete(ctx context.Context, id int64) error
}

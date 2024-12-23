package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
)

type HarvestRepository interface {
	Create(ctx context.Context, harvest *domain.Harvest) error
	FindAll(ctx context.Context) ([]*domain.Harvest, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error)
	FindByLandID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error)
	FindByLandCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error)
	FindByCityID(ctx context.Context, id int64) ([]*domain.Harvest, error)
	Update(ctx context.Context, id uuid.UUID, harvest *domain.Harvest) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindAllDeleted(ctx context.Context) ([]*domain.Harvest, error)
	FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error)
}

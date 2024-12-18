package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type SupplyRepository interface {
	Create(ctx context.Context, supply *domain.Supply) error
	FindAll(ctx context.Context) ([]*domain.Supply, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Supply, error)
	FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.Supply, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, supply *domain.Supply) error
	FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*domain.Supply, error)
}

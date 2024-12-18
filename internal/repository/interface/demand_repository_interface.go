package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type DemandRepository interface {
	Create(ctx context.Context, supply *domain.Demand) error
	FindAll(ctx context.Context) ([]*domain.Demand, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Demand, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Demand, error)
	FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.Demand, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, supply *domain.Demand) error
	FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*domain.Demand, error)
}

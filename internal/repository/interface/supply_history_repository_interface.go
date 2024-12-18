package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type SupplyHistoryRepository interface {
	Create(ctx context.Context, supply *domain.SupplyHistory) error
	FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) ([]*domain.SupplyHistory, error)
	// not used
	FindAll(ctx context.Context) ([]*domain.SupplyHistory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.SupplyHistory, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.SupplyHistory, error)
	FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.SupplyHistory, error)
}

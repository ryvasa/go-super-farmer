package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type DemandHistoryRepository interface {
	Create(ctx context.Context, supply *domain.DemandHistory) error
	FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) ([]*domain.DemandHistory, error)

	// not used
	FindAll(ctx context.Context) ([]*domain.DemandHistory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.DemandHistory, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.DemandHistory, error)
	FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.DemandHistory, error)
}

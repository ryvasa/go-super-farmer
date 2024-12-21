package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type SupplyHistoryRepository interface {
	Create(ctx context.Context, supply *domain.SupplyHistory) error
	FindByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.SupplyHistory, error)
	// not used
	FindAll(ctx context.Context) ([]*domain.SupplyHistory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.SupplyHistory, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.SupplyHistory, error)
	FindByCityID(ctx context.Context, id int64) ([]*domain.SupplyHistory, error)
}

package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
)

type DemandHistoryRepository interface {
	Create(ctx context.Context, supply *domain.DemandHistory) error
	FindByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.DemandHistory, error)

	// not used
	FindAll(ctx context.Context) ([]*domain.DemandHistory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.DemandHistory, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.DemandHistory, error)
	FindByCityID(ctx context.Context, id int64) ([]*domain.DemandHistory, error)
}

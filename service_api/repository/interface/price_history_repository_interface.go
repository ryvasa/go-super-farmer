package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
)

type PriceHistoryRepository interface {
	Create(ctx context.Context, priceHistory *domain.PriceHistory) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.PriceHistory, error)
	FindByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.PriceHistory, error)
}

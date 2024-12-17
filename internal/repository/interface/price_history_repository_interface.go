package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type PriceHistoryRepository interface {
	Create(ctx context.Context, priceHistory *domain.PriceHistory) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.PriceHistory, error)
	FindByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*[]domain.PriceHistory, error)
}

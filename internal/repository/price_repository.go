package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type PriceRepository interface {
	Create(ctx context.Context, price *domain.Price) error
	FindAll(ctx context.Context) (*[]domain.Price, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	FindByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Price, error)
	FindByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Price, error)
	Update(ctx context.Context, price *domain.Price) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
}

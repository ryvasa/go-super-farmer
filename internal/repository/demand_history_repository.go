package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type DemandHistoryRepository interface {
	Create(ctx context.Context, supply *domain.DemandHistory) error
	FindAll(ctx context.Context) (*[]domain.DemandHistory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.DemandHistory, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.DemandHistory, error)
	FindByRegionID(ctx context.Context, id uuid.UUID) (*[]domain.DemandHistory, error)
}

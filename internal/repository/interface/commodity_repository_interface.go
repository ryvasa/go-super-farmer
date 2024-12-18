package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/database/pagination"
)

type CommodityRepository interface {
	Create(ctx context.Context, land *domain.Commodity) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error)
	FindAll(ctx context.Context, params *pagination.PaginationParams) ([]domain.Commodity, error)
	Update(ctx context.Context, id uuid.UUID, land *domain.Commodity) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error)
}

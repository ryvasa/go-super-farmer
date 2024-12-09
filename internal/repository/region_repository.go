package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type RegionRepository interface {
	Create(ctx context.Context, region *domain.Region) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Region, error)
	FindAll(ctx context.Context) (*[]domain.Region, error)
	FindByProvinceID(ctx context.Context, id int64) (*[]domain.Region, error)
	Update(ctx context.Context, id uuid.UUID, region *domain.Region) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindDeleted(ctx context.Context, id uuid.UUID) (*domain.Region, error)
}

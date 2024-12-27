package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
)

type LandRepository interface {
	Create(ctx context.Context, land *domain.Land) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Land, error)
	FindByUserID(ctx context.Context, id uuid.UUID) ([]*domain.Land, error)
	FindAll(ctx context.Context) ([]*domain.Land, error)
	Update(ctx context.Context, id uuid.UUID, land *domain.Land) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Land, error)
	SumAllLandArea(ctx context.Context, params *dto.LandAreaParamsDTO) (float64, error)
}

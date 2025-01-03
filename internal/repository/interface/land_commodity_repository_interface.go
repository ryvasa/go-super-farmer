package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type LandCommodityRepository interface {
	Create(ctx context.Context, landCommodity *domain.LandCommodity) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error)
	FindByLandID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error)
	FindAll(ctx context.Context) ([]*domain.LandCommodity, error)
	FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error)
	Update(ctx context.Context, id uuid.UUID, landCommodity *domain.LandCommodity) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error)
	SumLandAreaByLandID(ctx context.Context, id uuid.UUID) (float64, error)
	SumNotHarvestedLandAreaByLandID(ctx context.Context, id uuid.UUID) (float64, error)
	SumAllLandCommodityArea(ctx context.Context, params *dto.LandAreaParamsDTO) (float64, error)
}

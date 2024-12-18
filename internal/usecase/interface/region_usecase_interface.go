package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type RegionUsecase interface {
	CreateRegion(ctx context.Context, req *dto.RegionCreateDto) (*domain.Region, error)
	GetAllRegions(ctx context.Context) ([]*domain.Region, error)
	GetRegionByID(ctx context.Context, id uuid.UUID) (*domain.Region, error)
}

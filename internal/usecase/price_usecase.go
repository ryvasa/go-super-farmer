package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type PriceUsecase interface {
	Create(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error)
	GetAll(ctx context.Context) (*[]dto.PriceResponseDTO, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	GetByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]dto.PriceResponseDTO, error)
	GetByRegionID(ctx context.Context, regionID uuid.UUID) (*[]dto.PriceResponseDTO, error)
	Update(ctx context.Context, req *dto.PriceUpdateDTO) (*dto.PriceResponseDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (*dto.PriceResponseDTO, error)
}

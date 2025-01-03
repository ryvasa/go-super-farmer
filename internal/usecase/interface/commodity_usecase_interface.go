package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type CommodityUsecase interface {
	CreateCommodity(ctx context.Context, req *dto.CommodityCreateDTO) (*domain.Commodity, error)
	GetAllCommodities(ctx context.Context, queryParams *dto.PaginationDTO) (*dto.PaginationResponseDTO, error)
	GetCommodityById(ctx context.Context, id uuid.UUID) (*domain.Commodity, error)
	UpdateCommodity(ctx context.Context, id uuid.UUID, req *dto.CommodityUpdateDTO) (*domain.Commodity, error)
	DeleteCommodity(ctx context.Context, id uuid.UUID) error
	RestoreCommodity(ctx context.Context, id uuid.UUID) (*domain.Commodity, error)
}

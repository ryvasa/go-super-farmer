package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type LandCommodityUsecase interface {
	CreateLandCommodity(ctx context.Context, req *dto.LandCommodityCreateDTO) (*domain.LandCommodity, error)
	GetLandCommodityByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error)
	GetLandCommodityByLandID(ctx context.Context, id uuid.UUID) (*[]domain.LandCommodity, error)
	GetLandCommodityByCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.LandCommodity, error)
	GetAllLandCommodity(ctx context.Context) (*[]domain.LandCommodity, error)
	UpdateLandCommodity(ctx context.Context, id uuid.UUID, req *dto.LandCommodityUpdateDTO) (*domain.LandCommodity, error)
	DeleteLandCommodity(ctx context.Context, id uuid.UUID) error
	RestoreLandCommodity(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error)
}

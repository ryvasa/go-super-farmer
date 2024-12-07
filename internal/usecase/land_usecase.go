package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type LandUsecase interface {
	CreateLand(ctx context.Context, userId uuid.UUID, req *dto.LandCreateDTO) (*domain.Land, error)
	GetLandByID(ctx context.Context, id uuid.UUID) (*domain.Land, error)
	GetLandByUserID(ctx context.Context, userID uuid.UUID) (*[]domain.Land, error)
	GetAllLands(ctx context.Context) (*[]domain.Land, error)
	UpdateLand(ctx context.Context, userId, id uuid.UUID, req *dto.LandUpdateDTO) (*domain.Land, error)
	DeleteLand(ctx context.Context, id uuid.UUID) error
	RestoreLand(ctx context.Context, id uuid.UUID) (*domain.Land, error)
}

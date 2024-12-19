package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type HarvestUsecase interface {
	CreateHarvest(ctx context.Context, req *dto.HarvestCreateDTO) (*domain.Harvest, error)
	GetAllHarvest(ctx context.Context) ([]*domain.Harvest, error)
	GetHarvestByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error)
	GetHarvestByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error)
	GetHarvestByLandID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error)
	GetHarvestByLandCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error)
	GetHarvestByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error)
	UpdateHarvest(ctx context.Context, id uuid.UUID, req *dto.HarvestUpdateDTO) (*domain.Harvest, error)
	DeleteHarvest(ctx context.Context, id uuid.UUID) error
	RestoreHarvest(ctx context.Context, id uuid.UUID) (*domain.Harvest, error)
	GetAllDeletedHarvest(ctx context.Context) ([]*domain.Harvest, error)
	GetHarvestDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error)
	DownloadHarvestByLandCommodityID(ctx context.Context, harvestParams *dto.HarvestParamsDTO) error
}

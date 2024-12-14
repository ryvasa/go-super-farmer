package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type DemandUsecase interface {
	CreateDemand(ctx context.Context, req *dto.DemandCreateDTO) (*domain.Demand, error)
	GetAllDemands(ctx context.Context) (*[]domain.Demand, error)
	GetDemandByID(ctx context.Context, id uuid.UUID) (*domain.Demand, error)
	GetDemandsByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Demand, error)
	GetDemandsByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Demand, error)
	UpdateDemand(ctx context.Context, id uuid.UUID, req *dto.DemandUpdateDTO) (*domain.Demand, error)
	DeleteDemand(ctx context.Context, id uuid.UUID) error
	GetDemandHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*[]domain.DemandHistory, error)
}

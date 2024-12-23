package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
)

type DemandUsecase interface {
	CreateDemand(ctx context.Context, req *dto.DemandCreateDTO) (*domain.Demand, error)
	GetAllDemands(ctx context.Context) ([]*domain.Demand, error)
	GetDemandByID(ctx context.Context, id uuid.UUID) (*domain.Demand, error)
	GetDemandsByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Demand, error)
	GetDemandsByCityID(ctx context.Context, cityID int64) ([]*domain.Demand, error)
	UpdateDemand(ctx context.Context, id uuid.UUID, req *dto.DemandUpdateDTO) (*domain.Demand, error)
	DeleteDemand(ctx context.Context, id uuid.UUID) error
	GetDemandHistoryByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.DemandHistory, error)
}

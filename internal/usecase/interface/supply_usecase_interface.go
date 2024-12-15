package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type SupplyUsecase interface {
	CreateSupply(ctx context.Context, req *dto.SupplyCreateDTO) (*domain.Supply, error)
	GetAllSupply(ctx context.Context) (*[]domain.Supply, error)
	GetSupplyByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error)
	GetSupplyByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Supply, error)
	GetSupplyByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Supply, error)
	UpdateSupply(ctx context.Context, id uuid.UUID, req *dto.SupplyUpdateDTO) (*domain.Supply, error)
	DeleteSupply(ctx context.Context, id uuid.UUID) error
	GetSupplyHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*[]domain.SupplyHistory, error)
}

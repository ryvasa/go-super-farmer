package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type PriceUsecase interface {
	CreatePrice(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error)
	GetAllPrices(ctx context.Context) (*[]domain.Price, error)
	GetPriceByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	GetPricesByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Price, error)
	GetPricesByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Price, error)
	UpdatePrice(ctx context.Context, id uuid.UUID, req *dto.PriceUpdateDTO) (*domain.Price, error)
	DeletePrice(ctx context.Context, id uuid.UUID) error
	RestorePrice(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	GetPriceByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*domain.Price, error)
	GetPriceHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*[]domain.PriceHistory, error)
	DownloadPriceHistoryByCommodityIDAndRegionID(ctx context.Context, params *dto.PriceParamsDTO) error
}

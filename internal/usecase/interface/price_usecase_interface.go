package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type PriceUsecase interface {
	CreatePrice(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error)
	GetAllPrices(ctx context.Context) ([]*domain.Price, error)
	GetPriceByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	GetPricesByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Price, error)
	GetPricesByCityID(ctx context.Context, cityID int64) ([]*domain.Price, error)
	UpdatePrice(ctx context.Context, id uuid.UUID, req *dto.PriceUpdateDTO) (*domain.Price, error)
	DeletePrice(ctx context.Context, id uuid.UUID) error
	RestorePrice(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	GetPriceByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) (*domain.Price, error)
	GetPriceHistoryByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.PriceHistory, error)
	DownloadPriceHistoryByCommodityIDAndCityID(ctx context.Context, params *dto.PriceParamsDTO) (*dto.DownloadResponseDTO, error)
	GetPriceExcelFile(ctx context.Context, params *dto.PriceParamsDTO) (*string, error)
}

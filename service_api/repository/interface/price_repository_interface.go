package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
)

type PriceRepository interface {
	Create(ctx context.Context, price *domain.Price) error
	FindAll(ctx context.Context, queryParams *dto.PaginationDTO) ([]*domain.Price, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	FindByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Price, error)
	FindByCityID(ctx context.Context, cityID int64) ([]*domain.Price, error)
	Update(ctx context.Context, id uuid.UUID, price *domain.Price) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	FindByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) (*domain.Price, error)
	Count(ctx context.Context, filter *dto.PaginationFilterDTO) (int64, error)
}

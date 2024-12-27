package repository_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
)

type SaleRepository interface {
	Create(ctx context.Context, sale *domain.Sale) error
	FindAll(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Sale, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error)
	FindByCommodityID(ctx context.Context, params *dto.PaginationDTO, id uuid.UUID) ([]*domain.Sale, error)
	FindByCityID(ctx context.Context, params *dto.PaginationDTO, id int64) ([]*domain.Sale, error)
	Update(ctx context.Context, id uuid.UUID, sale *domain.Sale) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	FindAllDeleted(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Sale, error)
	FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error)
	Count(ctx context.Context, filter *dto.ParamFilterDTO) (int64, error)
	DeletedCount(ctx context.Context, filter *dto.ParamFilterDTO) (int64, error)
}

package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type SaleUsecase interface {
	CreateSale(ctx context.Context, req *dto.SaleCreateDTO) (*domain.Sale, error)
	GetAllSales(ctx context.Context, pagination *dto.PaginationDTO) (*dto.PaginationResponseDTO, error)
	GetSaleByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error)
	GetSalesByCommodityID(ctx context.Context, pagination *dto.PaginationDTO, id uuid.UUID) (*dto.PaginationResponseDTO, error)
	GetSalesByCityID(ctx context.Context, pagination *dto.PaginationDTO, id int64) (*dto.PaginationResponseDTO, error)
	UpdateSale(ctx context.Context, id uuid.UUID, req *dto.SaleUpdateDTO) (*domain.Sale, error)
	DeleteSale(ctx context.Context, id uuid.UUID) error
	RestoreSale(ctx context.Context, id uuid.UUID) (*domain.Sale, error)
	GetAllDeletedSales(ctx context.Context, pagination *dto.PaginationDTO) (*dto.PaginationResponseDTO, error)
	GetDeletedSaleByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error)
}

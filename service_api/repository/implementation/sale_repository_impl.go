package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	"github.com/ryvasa/go-super-farmer/service_api/repository"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type SaleRepositoryImpl struct {
	repository.BaseRepository
}

func NewSaleRepository(db repository.BaseRepository) repository_interface.SaleRepository {
	return &SaleRepositoryImpl{db}
}

func (r *SaleRepositoryImpl) Create(ctx context.Context, sale *domain.Sale) error {
	return r.DB(ctx).Create(sale).Error
}

func (r *SaleRepositoryImpl) FindAll(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Sale, error) {
	var sales []*domain.Sale

	err := r.DB(ctx).
		Scopes(
			utils.ApplyFilters(&params.Filter),
			utils.GetPaginationScope(params),
		).
		Find(&sales).Error

	return sales, err
}

func (r *SaleRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	sale := domain.Sale{}

	err := r.DB(ctx).Where("id = ?", id).First(&sale).Error

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (r *SaleRepositoryImpl) FindByCommodityID(ctx context.Context, params *dto.PaginationDTO, id uuid.UUID) ([]*domain.Sale, error) {
	var sales []*domain.Sale

	err := r.DB(ctx).Scopes(
		utils.ApplyFilters(&params.Filter),
		utils.GetPaginationScope(params),
	).
		Where("commodity_id = ?", id).
		Find(&sales).Error

	return sales, err
}

func (r *SaleRepositoryImpl) FindByCityID(ctx context.Context, params *dto.PaginationDTO, id int64) ([]*domain.Sale, error) {
	var sales []*domain.Sale

	err := r.DB(ctx).Scopes(
		utils.ApplyFilters(&params.Filter),
		utils.GetPaginationScope(params),
	).Where("city_id = ?", id).Find(&sales).Error

	return sales, err
}

func (r *SaleRepositoryImpl) Update(ctx context.Context, id uuid.UUID, sale *domain.Sale) error {
	return r.DB(ctx).Model(&domain.Sale{}).Where("id = ?", id).Updates(sale).Error
}

func (r *SaleRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB(ctx).Delete(&domain.Sale{}, id).Error
}

func (r *SaleRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.DB(ctx).Unscoped().Model(&domain.Sale{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *SaleRepositoryImpl) FindAllDeleted(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Sale, error) {
	sales := []*domain.Sale{}

	err := r.DB(ctx).Unscoped().Scopes(
		utils.ApplyFilters(&params.Filter),
		utils.GetPaginationScope(params),
	).Where("deleted_at IS NOT NULL").Find(&sales).Error

	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (r *SaleRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	sale := domain.Sale{}

	err := r.DB(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&sale).Error

	if err != nil {
		return nil, err
	}

	return &sale, nil
}

func (r *SaleRepositoryImpl) Count(ctx context.Context, filter *dto.ParamFilterDTO) (int64, error) {
	var count int64
	err := r.DB(ctx).Model(&domain.Sale{}).
		Scopes(
			utils.ApplyFilters(filter),
		).Count(&count).Error
	return count, err
}

func (r *SaleRepositoryImpl) DeletedCount(ctx context.Context, filter *dto.ParamFilterDTO) (int64, error) {
	var count int64
	err := r.DB(ctx).Model(&domain.Sale{}).
		Scopes(
			utils.ApplyFilters(filter),
		).Where("deleted_at IS NOT NULL").
		Count(&count).Error
	return count, err
}

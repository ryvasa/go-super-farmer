package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"gorm.io/gorm"
)

type PriceRepositoryImpl struct {
	repository.BaseRepository
}

func NewPriceRepository(db repository.BaseRepository) repository_interface.PriceRepository {
	return &PriceRepositoryImpl{db}
}

func (r *PriceRepositoryImpl) Create(ctx context.Context, price *domain.Price) error {
	return r.DB(ctx).Create(price).Error
}

func (r *PriceRepositoryImpl) FindAll(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Price, error) {
	var prices []*domain.Price

	err := r.DB(ctx).
		Scopes(
			utils.ApplyFilters(&params.Filter),
			utils.GetPaginationScope(params),
		).
		Preload("Commodity").
		Preload("City").
		Preload("City.Province").
		Find(&prices).Error

	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (r *PriceRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	var price domain.Price

	err := r.DB(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("City").
		Preload("City.Province").
		First(&price, id).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PriceRepositoryImpl) FindByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Price, error) {
	var prices []*domain.Price

	err := r.DB(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("City").
		Preload("City.Province").
		Where("prices.commodity_id = ?", commodityID).
		Find(&prices).Error
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (r *PriceRepositoryImpl) FindByCityID(ctx context.Context, cityID int64) ([]*domain.Price, error) {
	var prices []*domain.Price
	err := r.DB(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("City").
		Preload("City.Province").
		Where("prices.city_id = ?", cityID).
		Find(&prices).Error
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func (r *PriceRepositoryImpl) Update(ctx context.Context, id uuid.UUID, price *domain.Price) error {
	return r.DB(ctx).Model(&domain.Price{}).Where("id = ?", id).Updates(price).Error

}

func (r *PriceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&domain.Price{}).Error

}

func (r *PriceRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.DB(ctx).Unscoped().Model(&domain.Price{}).Where("id = ?", id).Update("deleted_at", nil).Error

}

func (r *PriceRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	price := domain.Price{}
	err := r.DB(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("City").
		Preload("City.Province").
		Unscoped().
		Where("prices.id = ? AND prices.deleted_at IS NOT NULL", id).
		First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PriceRepositoryImpl) FindByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) (*domain.Price, error) {
	var price domain.Price
	err := r.DB(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("City").
		Preload("City.Province").
		Where("prices.commodity_id = ? AND prices.city_id = ?", commodityID, cityID).
		First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PriceRepositoryImpl) Count(ctx context.Context, filter *dto.ParamFilterDTO) (int64, error) {
	var count int64
	err := r.DB(ctx).
		Model(&domain.Price{}).
		Scopes(
			utils.ApplyFilters(filter),
		).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

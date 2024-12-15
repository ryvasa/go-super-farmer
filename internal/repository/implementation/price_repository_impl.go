package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type PriceRepositoryImpl struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) repository_interface.PriceRepository {
	return &PriceRepositoryImpl{db}
}

func (r *PriceRepositoryImpl) Create(ctx context.Context, price *domain.Price) error {
	err := r.db.WithContext(ctx).Create(price).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PriceRepositoryImpl) FindAll(ctx context.Context) (*[]domain.Price, error) {
	var prices []domain.Price

	err := r.db.WithContext(ctx).
		Preload("Commodity").
		Preload("Region").
		Preload("Region.Province").
		Preload("Region.City").
		Find(&prices).Error

	if err != nil {
		return nil, err
	}

	return &prices, nil
}

func (r *PriceRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	var price domain.Price
	err := r.db.WithContext(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("Region", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Region.Province").
		Preload("Region.City").
		First(&price, id).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PriceRepositoryImpl) FindByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Price, error) {
	var prices []domain.Price
	err := r.db.WithContext(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("Region", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Region.Province").
		Preload("Region.City").
		Where("prices.commodity_id = ?", commodityID).
		Find(&prices).Error
	if err != nil {
		return nil, err
	}
	return &prices, nil
}

func (r *PriceRepositoryImpl) FindByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Price, error) {
	var prices []domain.Price
	err := r.db.WithContext(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("Region", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Region.Province").
		Preload("Region.City").
		Where("prices.region_id = ?", regionID).
		Find(&prices).Error
	if err != nil {
		return nil, err
	}
	return &prices, nil
}

func (r *PriceRepositoryImpl) Update(ctx context.Context, id uuid.UUID, price *domain.Price) error {
	err := r.db.WithContext(ctx).Model(&domain.Price{}).Where("id = ?", id).Updates(price).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PriceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Price{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PriceRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Unscoped().Model(&domain.Price{}).Where("id = ?", id).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PriceRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	price := domain.Price{}
	err := r.db.WithContext(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("Region", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Region.Province").
		Preload("Region.City").
		Unscoped().
		Where("prices.id = ? AND prices.deleted_at IS NOT NULL", id).
		First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *PriceRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*domain.Price, error) {
	var price domain.Price
	err := r.db.WithContext(ctx).Model(&price).
		Where("prices.commodity_id = ? AND prices.region_id = ?", commodityID, regionID).
		First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

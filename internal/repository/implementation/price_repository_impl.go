package repository_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository/cache"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type PriceRepositoryImpl struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewPriceRepository(db *gorm.DB, cache cache.Cache) repository_interface.PriceRepository {
	return &PriceRepositoryImpl{db, cache}
}

func (r *PriceRepositoryImpl) Create(ctx context.Context, price *domain.Price) error {
	return r.db.WithContext(ctx).Create(price).Error
}

func (r *PriceRepositoryImpl) FindAll(ctx context.Context) (*[]domain.Price, error) {
	var prices []domain.Price

	key := fmt.Sprintf("price_%s", "all")
	cachedPrice, err := r.cache.Get(ctx, key)
	if err == nil && cachedPrice != nil {
		err := json.Unmarshal(cachedPrice, &prices)
		if err != nil {
			return nil, err
		}
		return &prices, nil
	}

	err = r.db.WithContext(ctx).
		Preload("Commodity").
		Preload("Region").
		Preload("Region.Province").
		Preload("Region.City").
		Find(&prices).Error

	if err != nil {
		return nil, err
	}

	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return nil, err
	}
	r.cache.Set(ctx, key, pricesJSON, 4*time.Minute)

	return &prices, nil
}

func (r *PriceRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	var price domain.Price

	key := fmt.Sprintf("price_%s", id)
	cachedPrice, err := r.cache.Get(ctx, key)
	if err == nil && cachedPrice != nil {
		err := json.Unmarshal(cachedPrice, &price)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}

	err = r.db.WithContext(ctx).
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

	priceJSON, err := json.Marshal(price)
	if err != nil {
		return nil, err
	}
	r.cache.Set(ctx, key, priceJSON, 4*time.Minute)

	return &price, nil
}

func (r *PriceRepositoryImpl) FindByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Price, error) {
	var prices []domain.Price

	key := fmt.Sprintf("price_com_%s", commodityID)
	cachedPrice, err := r.cache.Get(ctx, key)
	if err == nil && cachedPrice != nil {
		err := json.Unmarshal(cachedPrice, &prices)
		if err != nil {
			return nil, err
		}
		return &prices, nil
	}

	err = r.db.WithContext(ctx).
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

	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return nil, err
	}
	r.cache.Set(ctx, key, pricesJSON, 4*time.Minute)

	return &prices, nil
}

func (r *PriceRepositoryImpl) FindByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Price, error) {
	var prices []domain.Price
	key := fmt.Sprintf("price_reg_%s", regionID)
	cachedPrice, err := r.cache.Get(ctx, key)
	if err == nil && cachedPrice != nil {
		err := json.Unmarshal(cachedPrice, &prices)
		if err != nil {
			return nil, err
		}
		return &prices, nil
	}
	err = r.db.WithContext(ctx).
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

	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return nil, err
	}
	r.cache.Set(ctx, key, pricesJSON, 4*time.Minute)
	return &prices, nil
}

func (r *PriceRepositoryImpl) Update(ctx context.Context, id uuid.UUID, price *domain.Price) error {
	return r.db.WithContext(ctx).Model(&domain.Price{}).Where("id = ?", id).Updates(price).Error

}

func (r *PriceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Price{}).Error

}

func (r *PriceRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Model(&domain.Price{}).Where("id = ?", id).Update("deleted_at", nil).Error

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
	key := fmt.Sprintf("price_com_%s_reg_%s", commodityID, regionID)
	cachedPrice, err := r.cache.Get(ctx, key)
	if err == nil && cachedPrice != nil {
		err := json.Unmarshal(cachedPrice, &price)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}
	err = r.db.WithContext(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("Region", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Region.Province").
		Preload("Region.City").
		Where("prices.commodity_id = ? AND prices.region_id = ?", commodityID, regionID).
		First(&price).Error
	if err != nil {
		return nil, err
	}

	priceJSON, err := json.Marshal(price)
	if err != nil {
		return nil, err
	}
	r.cache.Set(ctx, key, priceJSON, 4*time.Minute)

	return &price, nil
}

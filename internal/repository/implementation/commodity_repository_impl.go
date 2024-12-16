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

type CommodityRepositoryImpl struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewCommodityRepository(db *gorm.DB, cache cache.Cache) repository_interface.CommodityRepository {
	return &CommodityRepositoryImpl{db, cache}
}

func (r *CommodityRepositoryImpl) Create(ctx context.Context, commodity *domain.Commodity) error {
	return r.db.WithContext(ctx).Create(commodity).Error
}

func (r *CommodityRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	var commodity domain.Commodity
	key := fmt.Sprintf("commodity_%s", id)
	cached, err := r.cache.Get(ctx, key)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &commodity)
		if err != nil {
			return nil, err
		}
		return &commodity, nil
	}
	err = r.db.WithContext(ctx).First(&commodity, id).Error
	if err != nil {
		return nil, err
	}

	commodityJSON, err := json.Marshal(commodity)
	if err != nil {
		return nil, err
	}

	r.cache.Set(ctx, key, commodityJSON, 4*time.Minute)

	return &commodity, nil
}

func (r *CommodityRepositoryImpl) FindAll(ctx context.Context) (*[]domain.Commodity, error) {
	var commodities []domain.Commodity
	key := fmt.Sprintf("commodity_%s", "all")
	cached, err := r.cache.Get(ctx, key)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &commodities)
		if err != nil {
			return nil, err
		}
		return &commodities, nil
	}

	if err := r.db.WithContext(ctx).Find(&commodities).Error; err != nil {
		return nil, err
	}

	commoditiesJSON, err := json.Marshal(commodities)
	if err != nil {
		return nil, err
	}

	r.cache.Set(ctx, key, commoditiesJSON, 4*time.Minute)

	return &commodities, nil
}

func (r *CommodityRepositoryImpl) Update(ctx context.Context, id uuid.UUID, commodity *domain.Commodity) error {
	return r.db.WithContext(ctx).Model(&domain.Commodity{}).Where("id = ?", id).Updates(commodity).Error
}

func (r *CommodityRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Commodity{}).Error
}

func (r *CommodityRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Model(&domain.Commodity{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *CommodityRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	var commodity domain.Commodity
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&commodity).Error; err != nil {
		return nil, err
	}
	return &commodity, nil
}

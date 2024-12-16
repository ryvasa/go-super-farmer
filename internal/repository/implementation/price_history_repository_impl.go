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

type PriceHistoryRepositoryImpl struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewPriceHistoryRepository(db *gorm.DB, cache cache.Cache) repository_interface.PriceHistoryRepository {
	return &PriceHistoryRepositoryImpl{db, cache}
}

func (r *PriceHistoryRepositoryImpl) Create(ctx context.Context, priceHistory *domain.PriceHistory) error {
	return r.db.WithContext(ctx).Create(priceHistory).Error
}

func (r *PriceHistoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.PriceHistory, error) {
	var priceHistory domain.PriceHistory
	err := r.db.WithContext(ctx).First(&priceHistory, id).Error
	if err != nil {
		return nil, err
	}
	return &priceHistory, nil
}

func (r *PriceHistoryRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*[]domain.PriceHistory, error) {
	userID := "ID"
	cacheKey := fmt.Sprintf("price_history_%s_%s_%s", userID, commodityID, regionID)
	cachedPriceHistory, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedPriceHistory != nil {
		var priceHistories []domain.PriceHistory
		err := json.Unmarshal(cachedPriceHistory, &priceHistories)
		if err != nil {
			return nil, err
		}
		return &priceHistories, nil
	}

	priceHistories := []domain.PriceHistory{}
	err = r.db.WithContext(ctx).Preload("Commodity", func(db *gorm.DB) *gorm.DB {
		return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
	}).
		Preload("Region", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Region.Province").
		Preload("Region.City").Where("commodity_id = ? AND region_id = ?", commodityID, regionID).Find(&priceHistories).Error
	if err != nil {
		return nil, err
	}

	userJSON, _ := json.Marshal(priceHistories)
	r.cache.Set(ctx, cacheKey, userJSON, 1*time.Minute)

	return &priceHistories, nil
}

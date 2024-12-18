package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type PriceHistoryRepositoryImpl struct {
	repository.BaseRepository
}

func NewPriceHistoryRepository(db repository.BaseRepository) repository_interface.PriceHistoryRepository {
	return &PriceHistoryRepositoryImpl{db}
}

func (r *PriceHistoryRepositoryImpl) Create(ctx context.Context, priceHistory *domain.PriceHistory) error {
	return r.DB(ctx).Create(priceHistory).Error
}

func (r *PriceHistoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.PriceHistory, error) {
	var priceHistory domain.PriceHistory
	err := r.DB(ctx).First(&priceHistory, id).Error
	if err != nil {
		return nil, err
	}
	return &priceHistory, nil
}

func (r *PriceHistoryRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) ([]*domain.PriceHistory, error) {
	priceHistories := []*domain.PriceHistory{}
	err := r.DB(ctx).Preload("Commodity", func(db *gorm.DB) *gorm.DB {
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

	return priceHistories, nil
}

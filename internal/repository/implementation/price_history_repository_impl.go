package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type PriceHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewPriceHistoryRepository(db *gorm.DB) repository_interface.PriceHistoryRepository {
	return &PriceHistoryRepositoryImpl{db}
}

func (r *PriceHistoryRepositoryImpl) Create(ctx context.Context, priceHistory *domain.PriceHistory) error {
	return r.db.WithContext(ctx).Create(priceHistory).Error
}

func (r *PriceHistoryRepositoryImpl) FindAll(ctx context.Context) (*[]domain.PriceHistory, error) {
	priceHistories := []domain.PriceHistory{}
	if err := r.db.WithContext(ctx).Find(&priceHistories).Error; err != nil {
		return nil, err
	}
	return &priceHistories, nil
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
	priceHistories := []domain.PriceHistory{}
	err := r.db.WithContext(ctx).Where("commodity_id = ? AND region_id = ?", commodityID, regionID).Find(&priceHistories).Error
	if err != nil {
		return nil, err
	}
	return &priceHistories, nil
}

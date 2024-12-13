package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type SupplyHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewSupplyHistoryRepository(db *gorm.DB) SupplyHistoryRepository {
	return &SupplyHistoryRepositoryImpl{db}
}

func (r *SupplyHistoryRepositoryImpl) Create(ctx context.Context, supply *domain.SupplyHistory) error {
	return r.db.WithContext(ctx).Create(supply).Error
}

func (r *SupplyHistoryRepositoryImpl) FindAll(ctx context.Context) (*[]domain.SupplyHistory, error) {
	var supplies []domain.SupplyHistory
	if err := r.db.WithContext(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *SupplyHistoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.SupplyHistory, error) {
	var supply domain.SupplyHistory
	err := r.db.WithContext(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *SupplyHistoryRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.SupplyHistory, error) {
	var supplies []domain.SupplyHistory
	if err := r.db.WithContext(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *SupplyHistoryRepositoryImpl) FindByRegionID(ctx context.Context, id uuid.UUID) (*[]domain.SupplyHistory, error) {
	var supplies []domain.SupplyHistory
	if err := r.db.WithContext(ctx).Where("region_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

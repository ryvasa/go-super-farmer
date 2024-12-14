package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type DemandHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewDemandHistoryRepository(db *gorm.DB) DemandHistoryRepository {
	return &DemandHistoryRepositoryImpl{db}
}

func (r *DemandHistoryRepositoryImpl) Create(ctx context.Context, supply *domain.DemandHistory) error {
	return r.db.WithContext(ctx).Create(supply).Error
}

func (r *DemandHistoryRepositoryImpl) FindAll(ctx context.Context) (*[]domain.DemandHistory, error) {
	var supplies []domain.DemandHistory
	if err := r.db.WithContext(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *DemandHistoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.DemandHistory, error) {
	var supply domain.DemandHistory
	err := r.db.WithContext(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *DemandHistoryRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.DemandHistory, error) {
	var supplies []domain.DemandHistory
	if err := r.db.WithContext(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *DemandHistoryRepositoryImpl) FindByRegionID(ctx context.Context, id uuid.UUID) (*[]domain.DemandHistory, error) {
	var supplies []domain.DemandHistory
	if err := r.db.WithContext(ctx).Where("region_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *DemandHistoryRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*[]domain.DemandHistory, error) {
	var supplies []domain.DemandHistory
	if err := r.db.WithContext(ctx).Where("commodity_id = ? AND region_id = ?", commodityID, regionID).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type DemandRepositoryImpl struct {
	db *gorm.DB
}

func NewDemandRepository(db *gorm.DB) DemandRepository {
	return &DemandRepositoryImpl{db}
}

func (r *DemandRepositoryImpl) Create(ctx context.Context, supply *domain.Demand) error {
	return r.db.WithContext(ctx).Create(supply).Error
}

func (r *DemandRepositoryImpl) FindAll(ctx context.Context) (*[]domain.Demand, error) {
	var supplies []domain.Demand
	if err := r.db.WithContext(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *DemandRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Demand, error) {
	var supply domain.Demand
	err := r.db.WithContext(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *DemandRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.Demand, error) {
	var supplies []domain.Demand
	if err := r.db.WithContext(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *DemandRepositoryImpl) FindByRegionID(ctx context.Context, id uuid.UUID) (*[]domain.Demand, error) {
	var supplies []domain.Demand
	if err := r.db.WithContext(ctx).Where("region_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *DemandRepositoryImpl) Update(ctx context.Context, id uuid.UUID, supply *domain.Demand) error {
	return r.db.WithContext(ctx).Model(&domain.Demand{}).Where("id = ?", id).Updates(supply).Error
}

func (r *DemandRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Demand{}).Error
}

func (r *DemandRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*domain.Demand, error) {
	var supply domain.Demand
	err := r.db.WithContext(ctx).Where("commodity_id = ? AND region_id = ?", commodityID, regionID).First(&supply).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

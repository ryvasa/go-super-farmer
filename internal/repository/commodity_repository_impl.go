package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type CommodityRepositoryImpl struct {
	db *gorm.DB
}

func NewCommodityRepository(db *gorm.DB) CommodityRepository {
	return &CommodityRepositoryImpl{db}
}

func (r *CommodityRepositoryImpl) Create(ctx context.Context, commodity *domain.Commodity) error {
	err := r.db.WithContext(ctx).Create(commodity).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CommodityRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	var commodity domain.Commodity
	err := r.db.WithContext(ctx).First(&commodity, id).Error
	if err != nil {
		return nil, err
	}
	return &commodity, nil
}

func (r *CommodityRepositoryImpl) FindAll(ctx context.Context) (*[]domain.Commodity, error) {
	var commodities []domain.Commodity
	if err := r.db.WithContext(ctx).Find(&commodities).Error; err != nil {
		return nil, err
	}
	return &commodities, nil
}

func (r *CommodityRepositoryImpl) Update(ctx context.Context, id uuid.UUID, commodity *domain.Commodity) error {
	err := r.db.WithContext(ctx).Model(&domain.Commodity{}).Where("id = ?", id).Updates(commodity).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CommodityRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Commodity{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CommodityRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Unscoped().Model(&domain.Commodity{}).Where("id = ?", id).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CommodityRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	var commodity domain.Commodity
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&commodity).Error; err != nil {
		return nil, err
	}
	return &commodity, nil
}

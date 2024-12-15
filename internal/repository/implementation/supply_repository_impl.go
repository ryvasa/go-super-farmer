package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type SupplyRepositoryImpl struct {
	db *gorm.DB
}

func NewSupplyRepository(db *gorm.DB) repository_interface.SupplyRepository {
	return &SupplyRepositoryImpl{db}
}

func (r *SupplyRepositoryImpl) Create(ctx context.Context, supply *domain.Supply) error {
	return r.db.WithContext(ctx).Create(supply).Error
}

func (r *SupplyRepositoryImpl) FindAll(ctx context.Context) (*[]domain.Supply, error) {
	var supplies []domain.Supply
	if err := r.db.WithContext(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *SupplyRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error) {
	var supply domain.Supply
	err := r.db.WithContext(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *SupplyRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.Supply, error) {
	var supplies []domain.Supply
	if err := r.db.WithContext(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *SupplyRepositoryImpl) FindByRegionID(ctx context.Context, id uuid.UUID) (*[]domain.Supply, error) {
	var supplies []domain.Supply
	if err := r.db.WithContext(ctx).Where("region_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

func (r *SupplyRepositoryImpl) Update(ctx context.Context, id uuid.UUID, supply *domain.Supply) error {
	return r.db.WithContext(ctx).Model(&domain.Supply{}).Where("id = ?", id).Updates(supply).Error
}

func (r *SupplyRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Supply{}).Error
}

func (r *SupplyRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*domain.Supply, error) {
	var supply domain.Supply
	err := r.db.WithContext(ctx).Where("commodity_id = ? AND region_id = ?", commodityID, regionID).First(&supply).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

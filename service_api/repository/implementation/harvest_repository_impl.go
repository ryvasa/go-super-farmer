package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"gorm.io/gorm"
)

type HarvestRepositoryImpl struct {
	db *gorm.DB
}

func NewHarvestRepository(db *gorm.DB) repository_interface.HarvestRepository {
	return &HarvestRepositoryImpl{db}
}

func (r *HarvestRepositoryImpl) Create(ctx context.Context, harvest *domain.Harvest) error {
	return r.db.WithContext(ctx).Create(harvest).Error
}

func (r *HarvestRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Harvest, error) {
	var harvests []*domain.Harvest

	if err := r.db.WithContext(ctx).Find(&harvests).Error; err != nil {
		return nil, err
	}

	return harvests, nil
}

func (r *HarvestRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	var harvest domain.Harvest

	err := r.db.WithContext(ctx).First(&harvest, id).Error
	if err != nil {
		return nil, err
	}
	return &harvest, nil
}

func (r *HarvestRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	var harvests []*domain.Harvest

	if err := r.db.WithContext(ctx).
		Preload("LandCommodity").
		Joins("JOIN land_commodities ON harvests.land_commodity_id = land_commodities.id").
		Where("land_commodities.commodity_id = ?", id).
		Find(&harvests).Error; err != nil {
		return nil, err
	}
	return harvests, nil
}

func (r *HarvestRepositoryImpl) FindByLandID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	var harvests []*domain.Harvest

	if err := r.db.WithContext(ctx).
		Preload("LandCommodity").
		Joins("JOIN land_commodities ON harvests.land_commodity_id = land_commodities.id").
		Where("land_commodities.land_id = ?", id).
		Find(&harvests).Error; err != nil {
		return nil, err
	}

	return harvests, nil
}

func (r *HarvestRepositoryImpl) FindByLandCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	var harvests []*domain.Harvest

	if err := r.db.WithContext(ctx).
		Where("land_commodity_id = ?", id).Find(&harvests).Error; err != nil {
		return nil, err
	}
	return harvests, nil
}

func (r *HarvestRepositoryImpl) FindByCityID(ctx context.Context, id int64) ([]*domain.Harvest, error) {
	var harvests []*domain.Harvest

	if err := r.db.WithContext(ctx).Where("city_id = ?", id).Find(&harvests).Error; err != nil {
		return nil, err
	}
	return harvests, nil
}

func (r *HarvestRepositoryImpl) Update(ctx context.Context, id uuid.UUID, harvest *domain.Harvest) error {
	return r.db.WithContext(ctx).Model(&domain.Harvest{}).Where("id = ?", id).Updates(harvest).Error
}

func (r *HarvestRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Harvest{}).Error
}

func (r *HarvestRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Model(&domain.Harvest{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *HarvestRepositoryImpl) FindAllDeleted(ctx context.Context) ([]*domain.Harvest, error) {
	var harvests []*domain.Harvest
	if err := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Find(&harvests).Error; err != nil {
		return nil, err
	}
	return harvests, nil
}

func (r *HarvestRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	var harvest domain.Harvest
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&harvest).Error; err != nil {
		return nil, err
	}
	return &harvest, nil
}

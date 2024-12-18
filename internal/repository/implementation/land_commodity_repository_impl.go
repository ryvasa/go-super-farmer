package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type LandCommodityRepositoryImpl struct {
	db *gorm.DB
}

func NewLandCommodityRepository(db *gorm.DB) repository_interface.LandCommodityRepository {
	return &LandCommodityRepositoryImpl{db}
}

func (r *LandCommodityRepositoryImpl) Create(ctx context.Context, landCommodity *domain.LandCommodity) error {
	return r.db.WithContext(ctx).Create(landCommodity).Error
}

func (r *LandCommodityRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error) {
	var landCommodity domain.LandCommodity

	err := r.db.WithContext(ctx).
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("Land", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		First(&landCommodity, id).Error
	if err != nil {
		return nil, err
	}
	return &landCommodity, nil
}

func (r *LandCommodityRepositoryImpl) FindByLandID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error) {
	var landCommodities []*domain.LandCommodity

	if err := r.db.WithContext(ctx).Where("land_id = ?", id).Find(&landCommodities).Error; err != nil {
		return nil, err
	}

	return landCommodities, nil
}

func (r *LandCommodityRepositoryImpl) FindAll(ctx context.Context) ([]*domain.LandCommodity, error) {
	var landCommodities []*domain.LandCommodity
	if err := r.db.WithContext(ctx).Find(&landCommodities).Error; err != nil {
		return nil, err
	}
	return landCommodities, nil
}

func (r *LandCommodityRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error) {
	var landCommodities []*domain.LandCommodity
	if err := r.db.WithContext(ctx).Where("commodity_id = ?", id).Find(&landCommodities).Error; err != nil {
		return nil, err
	}

	return landCommodities, nil
}

func (r *LandCommodityRepositoryImpl) Update(ctx context.Context, id uuid.UUID, landCommodity *domain.LandCommodity) error {
	err := r.db.WithContext(ctx).Model(&domain.LandCommodity{}).Where("id = ?", id).Updates(landCommodity).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *LandCommodityRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.LandCommodity{}).Error
}

func (r *LandCommodityRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Model(&domain.LandCommodity{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *LandCommodityRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error) {
	var landCommodity domain.LandCommodity
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&landCommodity).Error; err != nil {
		return nil, err
	}
	return &landCommodity, nil
}

func (r *LandCommodityRepositoryImpl) SumLandAreaByLandID(ctx context.Context, id uuid.UUID) (float64, error) {
	var landArea float64
	err := r.db.WithContext(ctx).
		Model(&domain.LandCommodity{}).
		Where("land_id = ?", id).
		Select("COALESCE(SUM(land_area), 0)").
		Scan(&landArea).
		Error
	if err != nil {
		return 0, err
	}
	return landArea, nil
}

func (r *LandCommodityRepositoryImpl) SumLandAreaByCommodityID(ctx context.Context, id uuid.UUID) (float64, error) {
	var landArea float64
	err := r.db.WithContext(ctx).Model(&domain.LandCommodity{}).Where("commodity_id = ?", id).Select("SUM(land_area)").Scan(&landArea).Error
	if err != nil {
		return 0, err
	}
	return landArea, nil
}

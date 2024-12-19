package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type RegionRepositoryImpl struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) repository_interface.RegionRepository {
	return &RegionRepositoryImpl{db}
}

func (r *RegionRepositoryImpl) Create(ctx context.Context, region *domain.Region) error {
	err := r.db.WithContext(ctx).Create(&region).Error
	return err
}

func (r *RegionRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Region, error) {
	var region domain.Region
	err := r.db.WithContext(ctx).First(&region, id).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *RegionRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Region, error) {
	var regions []*domain.Region
	err := r.db.WithContext(ctx).Find(&regions).Error
	if err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *RegionRepositoryImpl) FindByProvinceID(ctx context.Context, id int64) ([]*domain.Region, error) {
	var regions []*domain.Region
	err := r.db.WithContext(ctx).Where("province_id = ?", id).Find(&regions).Error
	if err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *RegionRepositoryImpl) Update(ctx context.Context, id uuid.UUID, region *domain.Region) error {
	err := r.db.WithContext(ctx).Model(&domain.Region{}).Where("id = ?", id).Updates(region).Error
	return err
}

func (r *RegionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&domain.Region{}, id).Error
	return err
}

func (r *RegionRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Unscoped().Model(&domain.Region{}).Where("id = ?", id).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RegionRepositoryImpl) FindDeleted(ctx context.Context, id uuid.UUID) (*domain.Region, error) {
	var region domain.Region
	err := r.db.WithContext(ctx).Unscoped().First(&region, id).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

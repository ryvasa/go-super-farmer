package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type LandRepositoryImpl struct {
	db *gorm.DB
}

func NewLandRepository(db *gorm.DB) repository_interface.LandRepository {
	return &LandRepositoryImpl{db}
}

func (r *LandRepositoryImpl) Create(ctx context.Context, land *domain.Land) error {
	err := r.db.WithContext(ctx).Create(land).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *LandRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	var land domain.Land
	err := r.db.WithContext(ctx).First(&land, id).Error
	if err != nil {
		return nil, err
	}
	return &land, nil
}

func (r *LandRepositoryImpl) FindByUserID(ctx context.Context, id uuid.UUID) ([]*domain.Land, error) {
	var lands []*domain.Land
	if err := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&lands).Error; err != nil {
		return nil, err
	}
	return lands, nil
}

func (r *LandRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Land, error) {
	var lands []*domain.Land
	if err := r.db.WithContext(ctx).Find(&lands).Error; err != nil {
		return nil, err
	}
	return lands, nil
}

func (r *LandRepositoryImpl) Update(ctx context.Context, id uuid.UUID, land *domain.Land) error {
	err := r.db.WithContext(ctx).Model(&domain.Land{}).Where("id = ?", id).Updates(land).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *LandRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Land{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *LandRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Unscoped().Model(&domain.Land{}).Where("id = ?", id).Update("deleted_at", nil).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *LandRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	var land domain.Land
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&land).Error; err != nil {
		return nil, err
	}
	return &land, nil
}

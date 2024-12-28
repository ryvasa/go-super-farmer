package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	"gorm.io/gorm"
)

func applyFilters(filter *dto.LandAreaParamsDTO) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.CityID != 0 {
			db = db.Where("city_id = ?", filter.CityID)
		}
		if filter.CommodityID != uuid.Nil {
			db = db.
				Joins("JOIN land_commodities ON land_commodities.land_id = lands.id").
				Where("land_commodities.commodity_id = ?", filter.CommodityID)
		}
		return db
	}
}

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

func (r *LandRepositoryImpl) SumAllLandArea(ctx context.Context, params *dto.LandAreaParamsDTO) (float64, error) {
	var landArea float64
	err := r.db.WithContext(ctx).
		Scopes(
			applyFilters(params),
		).
		Model(&domain.Land{}).
		Select("COALESCE(SUM(lands.land_area), 0)"). // Gunakan COALESCE
		Scan(&landArea).Error
	if err != nil {
		return 0, err
	}
	return landArea, nil
}

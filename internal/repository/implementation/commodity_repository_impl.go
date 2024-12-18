package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/pagination"
	"gorm.io/gorm"
)

type CommodityRepositoryImpl struct {
	db *gorm.DB
}

func NewCommodityRepository(db *gorm.DB) repository_interface.CommodityRepository {
	return &CommodityRepositoryImpl{db}
}

func (r *CommodityRepositoryImpl) Create(ctx context.Context, commodity *domain.Commodity) error {
	return r.db.WithContext(ctx).Create(commodity).Error
}

func (r *CommodityRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	var commodity domain.Commodity
	err := r.db.WithContext(ctx).First(&commodity, id).Error
	if err != nil {
		return nil, err
	}
	return &commodity, nil
}

func (r *CommodityRepositoryImpl) FindAll(ctx context.Context, params *pagination.PaginationParams) ([]domain.Commodity, error) {
	var commodities []domain.Commodity

	if err := r.db.WithContext(ctx).
		Scopes(pagination.Paginate(commodities, params, r.db, "name")).
		Find(&commodities).Error; err != nil {
		return nil, err
	}

	return commodities, nil
}

func (r *CommodityRepositoryImpl) Update(ctx context.Context, id uuid.UUID, commodity *domain.Commodity) error {
	return r.db.WithContext(ctx).Model(&domain.Commodity{}).Where("id = ?", id).Updates(commodity).Error
}

func (r *CommodityRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Commodity{}).Error
}

func (r *CommodityRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Model(&domain.Commodity{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *CommodityRepositoryImpl) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	var commodity domain.Commodity
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&commodity).Error; err != nil {
		return nil, err
	}
	return &commodity, nil
}

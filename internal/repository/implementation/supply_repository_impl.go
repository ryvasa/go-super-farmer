package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
)

type SupplyRepositoryImpl struct {
	repository.BaseRepository
}

func NewSupplyRepository(db repository.BaseRepository) repository_interface.SupplyRepository {
	return &SupplyRepositoryImpl{db}
}

func (r *SupplyRepositoryImpl) Create(ctx context.Context, supply *domain.Supply) error {
	return r.DB(ctx).Create(supply).Error
}

func (r *SupplyRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Supply, error) {
	var supplies []*domain.Supply
	if err := r.DB(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *SupplyRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error) {
	var supply domain.Supply
	err := r.DB(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *SupplyRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Supply, error) {
	var supplies []*domain.Supply
	if err := r.DB(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *SupplyRepositoryImpl) FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.Supply, error) {
	var supplies []*domain.Supply
	if err := r.DB(ctx).Where("region_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *SupplyRepositoryImpl) Update(ctx context.Context, id uuid.UUID, supply *domain.Supply) error {
	return r.DB(ctx).Model(&domain.Supply{}).Where("id = ?", id).Updates(supply).Error
}

func (r *SupplyRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&domain.Supply{}).Error
}

func (r *SupplyRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*domain.Supply, error) {
	var supply domain.Supply
	err := r.DB(ctx).Where("commodity_id = ? AND region_id = ?", commodityID, regionID).First(&supply).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

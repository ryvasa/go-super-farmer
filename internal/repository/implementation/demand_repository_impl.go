package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
)

type DemandRepositoryImpl struct {
	repository.BaseRepository
}

func NewDemandRepository(db repository.BaseRepository) repository_interface.DemandRepository {
	return &DemandRepositoryImpl{db}
}

func (r *DemandRepositoryImpl) Create(ctx context.Context, supply *domain.Demand) error {
	return r.DB(ctx).Create(supply).Error
}

func (r *DemandRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Demand, error) {
	var supplies []*domain.Demand
	if err := r.DB(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *DemandRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Demand, error) {
	var supply domain.Demand
	err := r.DB(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *DemandRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Demand, error) {
	var supplies []*domain.Demand
	if err := r.DB(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *DemandRepositoryImpl) FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.Demand, error) {
	var supplies []*domain.Demand
	if err := r.DB(ctx).Where("region_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *DemandRepositoryImpl) Update(ctx context.Context, id uuid.UUID, supply *domain.Demand) error {
	return r.DB(ctx).Model(&domain.Demand{}).Where("id = ?", id).Updates(supply).Error
}

func (r *DemandRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&domain.Demand{}).Error
}

func (r *DemandRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) (*domain.Demand, error) {
	var supply domain.Demand
	err := r.DB(ctx).Where("commodity_id = ? AND region_id = ?", commodityID, regionID).First(&supply).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
)

type DemandHistoryRepositoryImpl struct {
	repository.BaseRepository
}

func NewDemandHistoryRepository(db repository.BaseRepository) repository_interface.DemandHistoryRepository {
	return &DemandHistoryRepositoryImpl{db}
}

func (r *DemandHistoryRepositoryImpl) Create(ctx context.Context, supply *domain.DemandHistory) error {
	return r.DB(ctx).Create(supply).Error
}

func (r *DemandHistoryRepositoryImpl) FindAll(ctx context.Context) ([]*domain.DemandHistory, error) {
	var supplies []*domain.DemandHistory
	if err := r.DB(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *DemandHistoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.DemandHistory, error) {
	var supply domain.DemandHistory
	err := r.DB(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *DemandHistoryRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.DemandHistory, error) {
	var supplies []*domain.DemandHistory
	if err := r.DB(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *DemandHistoryRepositoryImpl) FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.DemandHistory, error) {
	var supplies []*domain.DemandHistory
	if err := r.DB(ctx).Where("region_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *DemandHistoryRepositoryImpl) FindByCommodityIDAndRegionID(ctx context.Context, commodityID uuid.UUID, regionID uuid.UUID) ([]*domain.DemandHistory, error) {
	var supplies []*domain.DemandHistory
	if err := r.DB(ctx).Where("commodity_id = ? AND region_id = ?", commodityID, regionID).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

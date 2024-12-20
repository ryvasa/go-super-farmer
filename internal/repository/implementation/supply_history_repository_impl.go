package repository_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
)

type SupplyHistoryRepositoryImpl struct {
	repository.BaseRepository
}

func NewSupplyHistoryRepository(db repository.BaseRepository) repository_interface.SupplyHistoryRepository {
	return &SupplyHistoryRepositoryImpl{db}
}

func (r *SupplyHistoryRepositoryImpl) Create(ctx context.Context, supply *domain.SupplyHistory) error {
	return r.DB(ctx).Create(supply).Error
}

func (r *SupplyHistoryRepositoryImpl) FindAll(ctx context.Context) ([]*domain.SupplyHistory, error) {
	var supplies []*domain.SupplyHistory
	if err := r.DB(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *SupplyHistoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.SupplyHistory, error) {
	var supply domain.SupplyHistory
	err := r.DB(ctx).First(&supply, id).Error
	if err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *SupplyHistoryRepositoryImpl) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.SupplyHistory, error) {
	var supplies []*domain.SupplyHistory
	if err := r.DB(ctx).Where("commodity_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *SupplyHistoryRepositoryImpl) FindByCityID(ctx context.Context, id int64) ([]*domain.SupplyHistory, error) {
	var supplies []*domain.SupplyHistory
	if err := r.DB(ctx).Where("city_id = ?", id).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *SupplyHistoryRepositoryImpl) FindByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.SupplyHistory, error) {
	var supplies []*domain.SupplyHistory
	if err := r.DB(ctx).Where("commodity_id = ? AND city_id = ?", commodityID, cityID).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

package repository

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type CityRepositoryImpl struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return &CityRepositoryImpl{db}
}

func (r *CityRepositoryImpl) Create(ctx context.Context, city *domain.City) error {
	return r.db.WithContext(ctx).Create(city).Error
}

func (r *CityRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.City, error) {
	var city domain.City
	err := r.db.WithContext(ctx).First(&city, id).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *CityRepositoryImpl) FindAll(ctx context.Context) (*[]domain.City, error) {
	var cities []domain.City
	if err := r.db.WithContext(ctx).Find(&cities).Error; err != nil {
		return nil, err
	}
	return &cities, nil
}

func (r *CityRepositoryImpl) Update(ctx context.Context, id int64, city *domain.City) error {
	return r.db.WithContext(ctx).Model(&domain.City{}).Where("id = ?", id).Updates(city).Error
}

func (r *CityRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.City{}).Error
}

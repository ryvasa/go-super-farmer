package repository_implementation

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	"gorm.io/gorm"
)

type ProvinceRepositoryImpl struct {
	db *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) repository_interface.ProvinceRepository {
	return &ProvinceRepositoryImpl{db}
}

func (r *ProvinceRepositoryImpl) Create(ctx context.Context, province *domain.Province) error {
	return r.db.WithContext(ctx).Create(province).Error
}

func (r *ProvinceRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.Province, error) {
	var province domain.Province
	err := r.db.WithContext(ctx).First(&province, id).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (r *ProvinceRepositoryImpl) FindAll(ctx context.Context) ([]*domain.Province, error) {
	var provinces []*domain.Province
	if err := r.db.WithContext(ctx).Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}

func (r *ProvinceRepositoryImpl) Update(ctx context.Context, id int64, province *domain.Province) error {
	return r.db.WithContext(ctx).Model(&domain.Province{}).Where("id = ?", id).Updates(province).Error
}

func (r *ProvinceRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Province{}).Error
}

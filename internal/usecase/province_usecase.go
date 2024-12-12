package usecase

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type ProvinceUsecase interface {
	CreateProvince(ctx context.Context, req *dto.ProvinceCreateDTO) (*domain.Province, error)
	GetAllProvinces(ctx context.Context) (*[]domain.Province, error)
	GetProvinceByID(ctx context.Context, id int64) (*domain.Province, error)
	UpdateProvince(ctx context.Context, id int64, req *dto.ProvinceUpdateDTO) (*domain.Province, error)
	DeleteProvince(ctx context.Context, id int64) error
}

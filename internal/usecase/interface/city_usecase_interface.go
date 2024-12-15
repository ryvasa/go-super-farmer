package usecase_interface

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type CityUsecase interface {
	CreateCity(ctx context.Context, req *dto.CityCreateDTO) (*domain.City, error)
	GetAllCities(ctx context.Context) (*[]domain.City, error)
	GetCityByID(ctx context.Context, id int64) (*domain.City, error)
	UpdateCity(ctx context.Context, id int64, req *dto.CityUpdateDTO) (*domain.City, error)
	DeleteCity(ctx context.Context, id int64) error
}

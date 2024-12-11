package usecase

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type CityUsecaseImpl struct {
	repo repository.CityRepository
}

func NewCityUsecase(repo repository.CityRepository) CityUsecase {
	return &CityUsecaseImpl{repo}
}

func (uc *CityUsecaseImpl) CreateCity(ctx context.Context, req *dto.CityCreateDTO) (*domain.City, error) {
	city := domain.City{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	city.Name = req.Name
	city.ProvinceID = req.ProvinceID

	err := uc.repo.Create(ctx, &city)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdCity, err := uc.repo.FindByID(ctx, city.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdCity, nil
}

func (uc *CityUsecaseImpl) GetAllCities(ctx context.Context) (*[]domain.City, error) {
	cities, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return cities, nil
}
func (uc *CityUsecaseImpl) GetCityById(ctx context.Context, id int64) (*domain.City, error) {
	city, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}
	return city, nil
}

func (uc *CityUsecaseImpl) UpdateCity(ctx context.Context, id int64, req *dto.CityUpdateDTO) (*domain.City, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	city, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}

	city.Name = req.Name
	city.ProvinceID = req.ProvinceID

	err = uc.repo.Update(ctx, id, city)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedCity, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedCity, nil
}

func (uc *CityUsecaseImpl) DeleteCity(ctx context.Context, id int64) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("city not found")
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	return nil
}

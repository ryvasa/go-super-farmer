package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type RegionUsecaseImpl struct {
	regionRepo   repository.RegionRepository
	cityRepo     repository.CityRepository
	provinceRepo repository.ProvinceRepository
}

func NewRegionUsecase(regionRepo repository.RegionRepository, cityRepo repository.CityRepository, provinceRepo repository.ProvinceRepository) RegionUsecase {
	return &RegionUsecaseImpl{regionRepo, cityRepo, provinceRepo}
}

func (uc *RegionUsecaseImpl) CreateRegion(ctx context.Context, req *dto.RegionCreateDto) (*domain.Region, error) {
	region := domain.Region{}

	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	_, err := uc.cityRepo.FindByID(ctx, req.CityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}

	_, err = uc.provinceRepo.FindByID(ctx, req.ProvinceID)
	if err != nil {
		return nil, utils.NewNotFoundError("province not found")
	}

	region.CityID = req.CityID
	region.ProvinceID = req.ProvinceID
	region.ID = uuid.New()

	err = uc.regionRepo.Create(ctx, &region)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	createdRegion, err := uc.regionRepo.FindByID(ctx, region.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return createdRegion, nil
}

func (uc *RegionUsecaseImpl) GetAllRegions(ctx context.Context) (*[]domain.Region, error) {
	regions, err := uc.regionRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return regions, nil
}

func (uc *RegionUsecaseImpl) GetRegionByID(ctx context.Context, id uuid.UUID) (*domain.Region, error) {
	region, err := uc.regionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return region, nil
}

package usecase

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type HarvestUsecaseImpl struct {
	harvestRepo       repository.HarvestRepository
	regionRepo        repository.RegionRepository
	landCommodityRepo repository.LandCommodityRepository
}

func NewHarvestUsecase(harvestRepo repository.HarvestRepository, regionRepo repository.RegionRepository, landCommodityRepo repository.LandCommodityRepository) HarvestUsecase {
	return &HarvestUsecaseImpl{harvestRepo, regionRepo, landCommodityRepo}
}

func (h *HarvestUsecaseImpl) CreateHarvest(ctx context.Context, req *dto.HarvestCreateDTO) (*domain.Harvest, error) {
	harvest := domain.Harvest{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	region, err := h.regionRepo.FindByID(ctx, req.RegionID)
	if err != nil {
		return nil, utils.NewNotFoundError("region not found")
	}
	commodityLand, err := h.landCommodityRepo.FindByID(ctx, req.LandCommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("land commodity not found")
	}

	parseDate, err := time.Parse("2006-01-02", req.HarvestDate)
	if err != nil {
		return nil, utils.NewBadRequestError("harvest date format is invalid")
	}

	harvest.RegionID = region.ID
	harvest.LandCommodityID = commodityLand.ID
	harvest.Quantity = req.Quantity
	harvest.Unit = req.Unit
	harvest.HarvestDate = parseDate
	harvest.ID = uuid.New()

	log.Println(harvest, "hahahah")

	err = h.harvestRepo.Create(ctx, &harvest)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdHarvest, err := h.harvestRepo.FindByID(ctx, harvest.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdHarvest, nil
}

func (h *HarvestUsecaseImpl) GetAllHarvest(ctx context.Context) (*[]domain.Harvest, error) {
	harvests, err := h.harvestRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (h *HarvestUsecaseImpl) GetHarvestByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	harvest, err := h.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("harvest not found")
	}
	return harvest, nil
}

func (h *HarvestUsecaseImpl) GetHarvestByCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.Harvest, error) {
	harvests, err := h.harvestRepo.FindByCommodityID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (h *HarvestUsecaseImpl) GetHarvestByLandID(ctx context.Context, id uuid.UUID) (*[]domain.Harvest, error) {
	harvests, err := h.harvestRepo.FindByLandID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (h *HarvestUsecaseImpl) GetHarvestByLandCommodityID(ctx context.Context, id uuid.UUID) (*[]domain.Harvest, error) {
	harvests, err := h.harvestRepo.FindByLandCommodityID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (h *HarvestUsecaseImpl) GetHarvestByRegionID(ctx context.Context, id uuid.UUID) (*[]domain.Harvest, error) {
	harvests, err := h.harvestRepo.FindByRegionID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (h *HarvestUsecaseImpl) UpdateHarvest(ctx context.Context, id uuid.UUID, req *dto.HarvestUpdateDTO) (*domain.Harvest, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	harvest, err := h.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("harvest not found")
	}

	if req.HarvestDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.HarvestDate)
		if err != nil {
			return nil, utils.NewValidationError(err)
		}
		harvest.HarvestDate = parsed
	}

	harvest.Quantity = req.Quantity
	harvest.Unit = req.Unit

	err = h.harvestRepo.Update(ctx, id, harvest)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedHarvest, err := h.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return updatedHarvest, nil
}

func (h *HarvestUsecaseImpl) DeleteHarvest(ctx context.Context, id uuid.UUID) error {
	_, err := h.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("harvest not found")
	}
	err = h.harvestRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (h *HarvestUsecaseImpl) RestoreHarvest(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	_, err := h.harvestRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted harvest not found")
	}
	err = h.harvestRepo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	restoredHarvest, err := h.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return restoredHarvest, nil
}

func (h *HarvestUsecaseImpl) GetAllDeletedHarvest(ctx context.Context) (*[]domain.Harvest, error) {
	harvests, err := h.harvestRepo.FindAllDeleted(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (h *HarvestUsecaseImpl) GetHarvestDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	harvest, err := h.harvestRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted harvest not found")
	}
	return harvest, nil
}

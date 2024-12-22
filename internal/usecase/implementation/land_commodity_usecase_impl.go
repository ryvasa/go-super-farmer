package usecase_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/utils"
)

type LandCommodityUsecaseImpl struct {
	landCommodityRepo repository_interface.LandCommodityRepository
	landRepo          repository_interface.LandRepository
	commodityRepo     repository_interface.CommodityRepository
	cache             cache.Cache
}

func NewLandCommodityUsecase(landCommodityRepo repository_interface.LandCommodityRepository, landRepo repository_interface.LandRepository, commodityRepo repository_interface.CommodityRepository, cache cache.Cache) usecase_interface.LandCommodityUsecase {
	return &LandCommodityUsecaseImpl{landCommodityRepo, landRepo, commodityRepo, cache}
}

func (u *LandCommodityUsecaseImpl) CreateLandCommodity(ctx context.Context, req *dto.LandCommodityCreateDTO) (*domain.LandCommodity, error) {
	landcommondity := domain.LandCommodity{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	commodity, err := u.commodityRepo.FindByID(ctx, req.CommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	land, err := u.landRepo.FindByID(ctx, req.LandID)
	if err != nil {
		return nil, utils.NewNotFoundError("land not found")
	}

	landArea, err := u.landCommodityRepo.SumLandAreaByLandID(ctx, land.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	if landArea+req.LandArea > land.LandArea {
		return nil, utils.NewBadRequestError("land area not enough")
	}

	landcommondity.LandID = land.ID
	landcommondity.CommodityID = commodity.ID
	landcommondity.LandArea = req.LandArea
	landcommondity.ID = uuid.New()

	err = u.landCommodityRepo.Create(ctx, &landcommondity)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	createdLandCommodity, err := u.landCommodityRepo.FindByID(ctx, landcommondity.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdLandCommodity, nil
}

func (u *LandCommodityUsecaseImpl) GetLandCommodityByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error) {
	landCommodity, err := u.landCommodityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("land commodity not found")
	}
	return landCommodity, nil
}

func (u *LandCommodityUsecaseImpl) GetLandCommodityByLandID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error) {
	landsCommodities, err := u.landCommodityRepo.FindByLandID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return landsCommodities, nil
}

func (u *LandCommodityUsecaseImpl) GetAllLandCommodity(ctx context.Context) ([]*domain.LandCommodity, error) {
	landCommodities := []*domain.LandCommodity{}
	key := fmt.Sprintf("land_commodity_%s", "all")
	cached, err := u.cache.Get(ctx, key)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &landCommodities)
		if err != nil {
			return nil, err
		}
		return landCommodities, nil
	}
	landCommodities, err = u.landCommodityRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())

	}

	landComJSON, err := json.Marshal(landCommodities)
	if err != nil {
		return nil, err
	}
	err = u.cache.Set(ctx, key, landComJSON, 4*time.Minute)

	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return landCommodities, nil
}

func (u *LandCommodityUsecaseImpl) GetLandCommodityByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error) {
	landCommodities, err := u.landCommodityRepo.FindByCommodityID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return landCommodities, nil
}

func (u *LandCommodityUsecaseImpl) UpdateLandCommodity(ctx context.Context, id uuid.UUID, req *dto.LandCommodityUpdateDTO) (*domain.LandCommodity, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	landCommodity, err := u.landCommodityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("land commodity not found")
	}

	_, err = u.commodityRepo.FindByID(ctx, req.CommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	land, err := u.landRepo.FindByID(ctx, req.LandID)
	if err != nil {
		return nil, utils.NewNotFoundError("land not found")
	}

	landArea, err := u.landCommodityRepo.SumLandAreaByLandID(ctx, land.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	// tidak tertes
	if (landArea-landCommodity.LandArea)+req.LandArea > land.LandArea {
		return nil, utils.NewBadRequestError("land area not enough")
	}

	landCommodity.LandArea = req.LandArea
	landCommodity.CommodityID = req.CommodityID
	landCommodity.LandID = req.LandID

	err = u.landCommodityRepo.Update(ctx, id, landCommodity)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedLandCommodity, err := u.landCommodityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedLandCommodity, nil
}

func (u *LandCommodityUsecaseImpl) DeleteLandCommodity(ctx context.Context, id uuid.UUID) error {
	_, err := u.landCommodityRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("land commodity not found")
	}
	err = u.landCommodityRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (u *LandCommodityUsecaseImpl) RestoreLandCommodity(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error) {
	deletedLandCommodity, err := u.landCommodityRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted land commodity not found")
	}

	land, err := u.landRepo.FindByID(ctx, deletedLandCommodity.LandID)
	if err != nil {
		return nil, utils.NewNotFoundError("land not found")
	}

	landArea, err := u.landCommodityRepo.SumLandAreaByLandID(ctx, land.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	if landArea+deletedLandCommodity.LandArea > land.LandArea {
		return nil, utils.NewBadRequestError("land area not enough")
	}

	err = u.landCommodityRepo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	restoredLandCommodity, err := u.landCommodityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return restoredLandCommodity, nil
}

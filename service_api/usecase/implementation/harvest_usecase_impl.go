package usecase_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/service_api/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/utils"
)

type HarvestMessage struct {
	LandCommodityID uuid.UUID `json:"LandCommodityID"`
	StartDate       time.Time `json:"StartDate"`
	EndDate         time.Time `json:"EndDate"`
}
type HarvestUsecaseImpl struct {
	harvestRepo       repository_interface.HarvestRepository
	cityRepo          repository_interface.CityRepository
	landCommodityRepo repository_interface.LandCommodityRepository
	rabbitMQ          messages.RabbitMQ
	cache             cache.Cache
	globFunc          utils.GlobFunc
}

func NewHarvestUsecase(harvestRepo repository_interface.HarvestRepository, cityRepo repository_interface.CityRepository, landCommodityRepo repository_interface.LandCommodityRepository, rabbitMQ messages.RabbitMQ, cache cache.Cache, globFunc utils.GlobFunc) usecase_interface.HarvestUsecase {
	return &HarvestUsecaseImpl{harvestRepo, cityRepo, landCommodityRepo, rabbitMQ, cache, globFunc}
}

func (uc *HarvestUsecaseImpl) CreateHarvest(ctx context.Context, req *dto.HarvestCreateDTO) (*domain.Harvest, error) {
	harvest := domain.Harvest{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	city, err := uc.cityRepo.FindByID(ctx, req.CityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}
	commodityLand, err := uc.landCommodityRepo.FindByID(ctx, req.LandCommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("land commodity not found")
	}

	parseDate, err := time.Parse("2006-01-02", req.HarvestDate)
	if err != nil {
		return nil, utils.NewBadRequestError("harvest date format is invalid")
	}

	harvest.CityID = city.ID
	harvest.LandCommodityID = commodityLand.ID
	harvest.Quantity = req.Quantity
	harvest.Unit = req.Unit
	harvest.HarvestDate = parseDate
	harvest.ID = uuid.New()

	err = uc.harvestRepo.Create(ctx, &harvest)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdHarvest, err := uc.harvestRepo.FindByID(ctx, harvest.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "harvest")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdHarvest, nil
}

func (uc *HarvestUsecaseImpl) GetAllHarvest(ctx context.Context) ([]*domain.Harvest, error) {
	var harvests []*domain.Harvest
	key := fmt.Sprintf("harvest_%s", "all")
	cached, err := uc.cache.Get(ctx, key)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &harvests)
		if err != nil {
			return nil, err
		}
		return harvests, nil
	}
	harvests, err = uc.harvestRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	harvestsJSON, err := json.Marshal(harvests)
	if err != nil {
		return nil, err
	}
	err = uc.cache.Set(ctx, key, harvestsJSON, 4*time.Minute)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return harvests, nil
}

func (uc *HarvestUsecaseImpl) GetHarvestByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	harvest, err := uc.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("harvest not found")
	}
	return harvest, nil
}

func (uc *HarvestUsecaseImpl) GetHarvestByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	harvests, err := uc.harvestRepo.FindByCommodityID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (uc *HarvestUsecaseImpl) GetHarvestByLandID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	harvests, err := uc.harvestRepo.FindByLandID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (uc *HarvestUsecaseImpl) GetHarvestByLandCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	harvests, err := uc.harvestRepo.FindByLandCommodityID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (uc *HarvestUsecaseImpl) GetHarvestByCityID(ctx context.Context, id int64) ([]*domain.Harvest, error) {
	harvests, err := uc.harvestRepo.FindByCityID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (uc *HarvestUsecaseImpl) UpdateHarvest(ctx context.Context, id uuid.UUID, req *dto.HarvestUpdateDTO) (*domain.Harvest, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	harvest, err := uc.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("harvest not found")
	}

	if req.HarvestDate != "" {
		parsed, err := time.Parse("2006-01-02", req.HarvestDate)
		if err != nil {
			return nil, utils.NewValidationError(err)
		}
		harvest.HarvestDate = parsed
	}

	harvest.Quantity = req.Quantity
	harvest.Unit = req.Unit

	err = uc.harvestRepo.Update(ctx, id, harvest)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedHarvest, err := uc.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	err = uc.cache.DeleteByPattern(ctx, "harvest")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedHarvest, nil
}

func (uc *HarvestUsecaseImpl) DeleteHarvest(ctx context.Context, id uuid.UUID) error {
	_, err := uc.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("harvest not found")
	}
	err = uc.harvestRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	err = uc.cache.DeleteByPattern(ctx, "harvest")
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (uc *HarvestUsecaseImpl) RestoreHarvest(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	_, err := uc.harvestRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted harvest not found")
	}
	err = uc.harvestRepo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	restoredHarvest, err := uc.harvestRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "harvest")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return restoredHarvest, nil
}

func (uc *HarvestUsecaseImpl) GetAllDeletedHarvest(ctx context.Context) ([]*domain.Harvest, error) {
	harvests, err := uc.harvestRepo.FindAllDeleted(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return harvests, nil
}

func (uc *HarvestUsecaseImpl) GetHarvestDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	harvest, err := uc.harvestRepo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("deleted harvest not found")
	}
	return harvest, nil
}

func (uc *HarvestUsecaseImpl) DownloadHarvestByLandCommodityID(ctx context.Context, harvestParams *dto.HarvestParamsDTO) (*dto.DownloadResponseDTO, error) {

	msg := HarvestMessage{
		LandCommodityID: harvestParams.LandCommodityID,
		StartDate:       harvestParams.StartDate,
		EndDate:         harvestParams.EndDate,
	}

	err := uc.rabbitMQ.PublishJSON(ctx, "report-exchange", "harvest", msg)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	res := &dto.DownloadResponseDTO{
		Message: "Report generation in progress. Please check back in a few moments.",
		DownloadURL: fmt.Sprintf("http://localhost:8080/api/harvests/land_commodity/%s/download/file?start_date=%s&end_date=%s",
			harvestParams.LandCommodityID, harvestParams.StartDate.Format("2006-01-02"), harvestParams.EndDate.Format("2006-01-02")),
	}

	return res, nil
}

func (uc *HarvestUsecaseImpl) GetHarvestExcelFile(ctx context.Context, harvestParams *dto.HarvestParamsDTO) (*string, error) {
	// Get the latest excel file
	filePath := fmt.Sprintf("./public/reports/harvests_%s_%s_%s_*.xlsx", harvestParams.LandCommodityID, harvestParams.StartDate.Format("2006-01-02"), harvestParams.EndDate.Format("2006-01-02"))
	matches, err := uc.globFunc.Glob(filePath) // Gunakan globFunc yang bisa dimock
	if err != nil {
		return nil, utils.NewInternalError("Error finding report file")
	}

	if len(matches) == 0 {
		return nil, utils.NewNotFoundError("Report file not found")
	}

	// Get the latest file (assuming filename contains timestamp)
	latestFile := matches[len(matches)-1]

	return &latestFile, nil
}
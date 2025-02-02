package handler_implementation

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type HarvestHandlerImpl struct {
	uc usecase_interface.HarvestUsecase
}

func NewHarvestHandler(uc usecase_interface.HarvestUsecase) handler_interface.HarvestHandler {
	return &HarvestHandlerImpl{uc}
}

func (h *HarvestHandlerImpl) CreateHarvest(c *gin.Context) {
	var req dto.HarvestCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvest, err := h.uc.CreateHarvest(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, harvest)
}

func (h *HarvestHandlerImpl) GetAllHarvest(c *gin.Context) {
	harvests, err := h.uc.GetAllHarvest(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvest, err := h.uc.GetHarvestByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvest)
}

func (h *HarvestHandlerImpl) GetHarvestByCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvests, err := h.uc.GetHarvestByCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByLandID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvests, err := h.uc.GetHarvestByLandID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByLandCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	harvests, err := h.uc.GetHarvestByLandCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByCityID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvests, err := h.uc.GetHarvestByCityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) UpdateHarvest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	var req dto.HarvestUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedHarvest, err := h.uc.UpdateHarvest(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedHarvest)
}

func (h *HarvestHandlerImpl) DeleteHarvest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	err = h.uc.DeleteHarvest(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Harvest deleted successfully"})
}

func (h *HarvestHandlerImpl) RestoreHarvest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	restoredHarvest, err := h.uc.RestoreHarvest(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredHarvest)
}

func (h *HarvestHandlerImpl) GetAllDeletedHarvest(c *gin.Context) {
	harvests, err := h.uc.GetAllDeletedHarvest(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestDeletedByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvest, err := h.uc.GetHarvestDeletedByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvest)
}

func (h *HarvestHandlerImpl) DownloadHarvestByLandCommodityID(c *gin.Context) {
	harvestParams := &dto.HarvestParamsDTO{}

	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
			return
		}
		harvestParams.StartDate = startDate
	}
	endDatestr := c.Query("end_date")
	if endDatestr != "" {
		endDate, err := time.Parse("2006-01-02", endDatestr)
		if err != nil {
			utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
			return
		}
		harvestParams.EndDate = endDate
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvestParams.LandCommodityID = id
	res, err := h.uc.DownloadHarvestByLandCommodityID(c, harvestParams)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

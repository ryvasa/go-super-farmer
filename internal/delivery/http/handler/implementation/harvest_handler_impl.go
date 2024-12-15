package handler_implementation

import (
	"net/http"

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

func (h *HarvestHandlerImpl) GetHarvestByRegionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvests, err := h.uc.GetHarvestByRegionID(c, id)
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

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type DemandHandlerImpl struct {
	demandUsecase usecase.DemandUsecase
}

func NewDemandHandler(demandUsecase usecase.DemandUsecase) DemandHandler {
	return &DemandHandlerImpl{demandUsecase}
}

func (h *DemandHandlerImpl) CreateDemand(c *gin.Context) {
	var req dto.DemandCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	demand, err := h.demandUsecase.CreateDemand(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, demand)
}

func (h *DemandHandlerImpl) GetAllDemands(c *gin.Context) {
	demands, err := h.demandUsecase.GetAllDemands(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, demands)
}

func (h *DemandHandlerImpl) GetDemandByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	demand, err := h.demandUsecase.GetDemandByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, demand)
}

func (h *DemandHandlerImpl) GetDemandsByCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	demands, err := h.demandUsecase.GetDemandsByCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, demands)
}

func (h *DemandHandlerImpl) GetDemandsByRegionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	demands, err := h.demandUsecase.GetDemandsByRegionID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, demands)
}

func (h *DemandHandlerImpl) UpdateDemand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	var req dto.DemandUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedDemand, err := h.demandUsecase.UpdateDemand(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedDemand)
}
func (h *DemandHandlerImpl) DeleteDemand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	err = h.demandUsecase.DeleteDemand(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Demand deleted successfully"})
}

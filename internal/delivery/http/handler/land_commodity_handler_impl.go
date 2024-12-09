package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type LandCommodityHandlerImpl struct {
	landCommodityUsecase usecase.LandCommodityUsecase
}

func NewLandCommodityHandler(landCommodityUsecase usecase.LandCommodityUsecase) LandCommodityHandler {
	return &LandCommodityHandlerImpl{landCommodityUsecase}
}

func (h *LandCommodityHandlerImpl) CreateLandCommodity(c *gin.Context) {
	var req dto.LandCommodityCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	landCommodity, err := h.landCommodityUsecase.CreateLandCommodity(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, landCommodity)
}

func (h *LandCommodityHandlerImpl) GetLandCommodityByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	landCommodity, err := h.landCommodityUsecase.GetLandCommodityByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, landCommodity)
}

func (h *LandCommodityHandlerImpl) GetLandCommodityByLandID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	landsCommodities, err := h.landCommodityUsecase.GetLandCommodityByLandID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, landsCommodities)
}

func (h *LandCommodityHandlerImpl) GetAllLandCommodity(c *gin.Context) {
	landCommodities, err := h.landCommodityUsecase.GetAllLandCommodity(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, landCommodities)
}

func (h *LandCommodityHandlerImpl) GetLandCommodityByCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	landCommodities, err := h.landCommodityUsecase.GetLandCommodityByCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, landCommodities)
}

func (h *LandCommodityHandlerImpl) UpdateLandCommodity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	var req dto.LandCommodityUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedLandCommodity, err := h.landCommodityUsecase.UpdateLandCommodity(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedLandCommodity)
}

func (h *LandCommodityHandlerImpl) DeleteLandCommodity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	err = h.landCommodityUsecase.DeleteLandCommodity(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Land commodity deleted successfully"})
}

func (h *LandCommodityHandlerImpl) RestoreLandCommodity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	restoredLandCommodity, err := h.landCommodityUsecase.RestoreLandCommodity(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredLandCommodity)
}

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

type LandCommodityHandlerImpl struct {
	uc usecase_interface.LandCommodityUsecase
}

func NewLandCommodityHandler(uc usecase_interface.LandCommodityUsecase) handler_interface.LandCommodityHandler {
	return &LandCommodityHandlerImpl{uc}
}

func (h *LandCommodityHandlerImpl) CreateLandCommodity(c *gin.Context) {
	var req dto.LandCommodityCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	landCommodity, err := h.uc.CreateLandCommodity(c, &req)
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

	landCommodity, err := h.uc.GetLandCommodityByID(c, id)
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
	landsCommodities, err := h.uc.GetLandCommodityByLandID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, landsCommodities)
}

func (h *LandCommodityHandlerImpl) GetAllLandCommodity(c *gin.Context) {
	landCommodities, err := h.uc.GetAllLandCommodity(c)
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
	landCommodities, err := h.uc.GetLandCommodityByCommodityID(c, id)
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
	updatedLandCommodity, err := h.uc.UpdateLandCommodity(c, id, &req)
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

	err = h.uc.DeleteLandCommodity(c, id)
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
	restoredLandCommodity, err := h.uc.RestoreLandCommodity(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredLandCommodity)
}

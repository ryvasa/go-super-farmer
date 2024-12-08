package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type CommodityHandlerImpl struct {
	commodityUsecase usecase.CommodityUsecase
}

func NewCommodityHandler(commodityUsecase usecase.CommodityUsecase) CommodityHandler {
	return &CommodityHandlerImpl{commodityUsecase}
}

func (h *CommodityHandlerImpl) CreateCommodity(c *gin.Context) {
	var req dto.CommodityCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	commodity, err := h.commodityUsecase.CreateCommodity(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, commodity)
}

func (h *CommodityHandlerImpl) GetAllCommodities(c *gin.Context) {
	commodities, err := h.commodityUsecase.GetAllCommodities(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, commodities)
}

func (h *CommodityHandlerImpl) GetCommodityById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	commodity, err := h.commodityUsecase.GetCommodityById(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	if commodity == nil {
		utils.ErrorResponse(c, utils.NewNotFoundError("commodity not found"))
		return
	}
	utils.SuccessResponse(c, http.StatusOK, commodity)
}

func (h *CommodityHandlerImpl) UpdateCommodity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	var req dto.CommodityUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedCommodity, err := h.commodityUsecase.UpdateCommodity(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedCommodity)
}

func (h *CommodityHandlerImpl) DeleteCommodity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	if err := h.commodityUsecase.DeleteCommodity(c, id); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Commodity deleted successfully"})
}

func (h *CommodityHandlerImpl) RestoreCommodity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	restoredCommodity, err := h.commodityUsecase.RestoreCommodity(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredCommodity)
}

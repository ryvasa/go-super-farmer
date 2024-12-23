package handler_implementation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type CommodityHandlerImpl struct {
	uc usecase_interface.CommodityUsecase
}

func NewCommodityHandler(uc usecase_interface.CommodityUsecase) handler_interface.CommodityHandler {
	return &CommodityHandlerImpl{uc}
}

func (h *CommodityHandlerImpl) CreateCommodity(c *gin.Context) {
	var req dto.CommodityCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	commodity, err := h.uc.CreateCommodity(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, commodity)
}

func (h *CommodityHandlerImpl) GetAllCommodities(c *gin.Context) {
	pagination, err := utils.GetPaginationParams(c)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	response, err := h.uc.GetAllCommodities(c, pagination)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, response)
}

func (h *CommodityHandlerImpl) GetCommodityById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	commodity, err := h.uc.GetCommodityById(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
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
	updatedCommodity, err := h.uc.UpdateCommodity(c, id, &req)
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
	if err := h.uc.DeleteCommodity(c, id); err != nil {
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
	restoredCommodity, err := h.uc.RestoreCommodity(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredCommodity)
}

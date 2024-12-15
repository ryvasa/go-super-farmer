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

type SupplyHandlerImpl struct {
	uc usecase_interface.SupplyUsecase
}

func NewSupplyHandler(uc usecase_interface.SupplyUsecase) handler_interface.SupplyHandler {
	return &SupplyHandlerImpl{uc: uc}
}

func (h *SupplyHandlerImpl) CreateSupply(c *gin.Context) {
	var req dto.SupplyCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	supply, err := h.uc.CreateSupply(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, supply)
}

func (h *SupplyHandlerImpl) GetAllSupply(c *gin.Context) {
	supplies, err := h.uc.GetAllSupply(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, supplies)
}

func (h *SupplyHandlerImpl) GetSupplyByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	supply, err := h.uc.GetSupplyByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, supply)
}

func (h *SupplyHandlerImpl) GetSupplyByCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	supplies, err := h.uc.GetSupplyByCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, supplies)
}

func (h *SupplyHandlerImpl) GetSupplyByRegionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	supplies, err := h.uc.GetSupplyByRegionID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, supplies)
}

func (h *SupplyHandlerImpl) UpdateSupply(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	var req dto.SupplyUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedSupply, err := h.uc.UpdateSupply(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedSupply)
}

func (h *SupplyHandlerImpl) DeleteSupply(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	err = h.uc.DeleteSupply(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Supply deleted successfully"})
}

func (h *SupplyHandlerImpl) GetSupplyHistoryByCommodityIDAndRegionID(c *gin.Context) {
	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	regionID, err := uuid.Parse(c.Param("region_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	supplies, err := h.uc.GetSupplyHistoryByCommodityIDAndRegionID(c, commodityID, regionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, supplies)
}

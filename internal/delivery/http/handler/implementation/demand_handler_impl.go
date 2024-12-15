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

type DemandHandlerImpl struct {
	uc usecase_interface.DemandUsecase
}

func NewDemandHandler(uc usecase_interface.DemandUsecase) handler_interface.DemandHandler {
	return &DemandHandlerImpl{uc}
}

func (h *DemandHandlerImpl) CreateDemand(c *gin.Context) {
	var req dto.DemandCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	demand, err := h.uc.CreateDemand(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, demand)
}

func (h *DemandHandlerImpl) GetAllDemands(c *gin.Context) {
	demands, err := h.uc.GetAllDemands(c)
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

	demand, err := h.uc.GetDemandByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, demand)
}

func (h *DemandHandlerImpl) GetDemandsByCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	demands, err := h.uc.GetDemandsByCommodityID(c, id)
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
	demands, err := h.uc.GetDemandsByRegionID(c, id)
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
	updatedDemand, err := h.uc.UpdateDemand(c, id, &req)
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

	err = h.uc.DeleteDemand(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Demand deleted successfully"})
}

func (h *DemandHandlerImpl) GetDemandHistoryByCommodityIDAndRegionID(c *gin.Context) {
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
	demands, err := h.uc.GetDemandHistoryByCommodityIDAndRegionID(c, commodityID, regionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, demands)
}

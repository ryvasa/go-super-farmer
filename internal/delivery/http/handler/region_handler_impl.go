package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type RegionHandlerImpl struct {
	uc usecase.RegionUseCase
}

func NewRegionHandler(uc usecase.RegionUseCase) RegionHandler {
	return &RegionHandlerImpl{uc}
}

func (h *RegionHandlerImpl) CreateRegion(c *gin.Context) {
	req := dto.RegionCreateDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	region, err := h.uc.CreateRegion(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, region)
}

func (h *RegionHandlerImpl) GetAllRegions(c *gin.Context) {
	regions, err := h.uc.GetAllRegions(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, regions)
}

func (h *RegionHandlerImpl) GetRegionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	region, err := h.uc.GetRegionByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, region)
}

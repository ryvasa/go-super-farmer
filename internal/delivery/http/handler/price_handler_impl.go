package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceHandlerImpl struct {
	uc usecase.PriceUsecase
}

func NewPriceHandler(uc usecase.PriceUsecase) PriceHandler {
	return &PriceHandlerImpl{uc}
}

func (h *PriceHandlerImpl) CreatePrice(c *gin.Context) {
	var req dto.PriceCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	createdPrice, err := h.uc.Create(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, createdPrice)
}

func (h *PriceHandlerImpl) GetAllPrices(c *gin.Context) {
	// TODO: Implement
}

func (h *PriceHandlerImpl) GetPriceById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	user, err := h.uc.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, user)
}
func (h *PriceHandlerImpl) GetPricesByCommodityID(c *gin.Context) {
	// TODO: Implement
}

func (h *PriceHandlerImpl) GetPricesByRegionID(c *gin.Context) {
	// TODO: Implement
}

func (h *PriceHandlerImpl) UpdatePrice(c *gin.Context) {
	// TODO: Implement
}

func (h *PriceHandlerImpl) DeletePrice(c *gin.Context) {
	// TODO: Implement
}

func (h *PriceHandlerImpl) RestorePrice(c *gin.Context) {
	// TODO: Implement
}

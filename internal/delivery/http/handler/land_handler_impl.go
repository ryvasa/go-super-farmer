package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type LandHandlerImpl struct {
	usecase usecase.LandUsecase
	// authUtil utils.AuthUtil
}

func NewLandHandler(usecase usecase.LandUsecase) LandHandler {
	return &LandHandlerImpl{usecase}
}

func (h *LandHandlerImpl) CreateLand(c *gin.Context) {
	userId, err := utils.GetAuthUserID(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	var req dto.LandCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	land, err := h.usecase.CreateLand(c, userId, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, land)
}

func (h *LandHandlerImpl) GetLandByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	land, err := h.usecase.GetLandByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	if land == nil {
		utils.ErrorResponse(c, utils.NewNotFoundError("land not found"))
		return
	}
	utils.SuccessResponse(c, http.StatusOK, land)
}

func (h *LandHandlerImpl) GetLandByUserID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	lands, err := h.usecase.GetLandByUserID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, lands)
}

func (h *LandHandlerImpl) GetAllLands(c *gin.Context) {
	lands, err := h.usecase.GetAllLands(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, lands)
}

func (h *LandHandlerImpl) UpdateLand(c *gin.Context) {
	userId, err := utils.GetAuthUserID(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	var req dto.LandUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedLand, err := h.usecase.UpdateLand(c, userId, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedLand)
}

func (h *LandHandlerImpl) DeleteLand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	err = h.usecase.DeleteLand(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Land deleted successfully"})
}

func (h *LandHandlerImpl) RestoreLand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	restoredLand, err := h.usecase.RestoreLand(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredLand)

}

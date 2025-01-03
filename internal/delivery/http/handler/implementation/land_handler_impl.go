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

type LandHandlerImpl struct {
	uc       usecase_interface.LandUsecase
	authUtil utils.AuthUtil
}

func NewLandHandler(uc usecase_interface.LandUsecase, authUtil utils.AuthUtil) handler_interface.LandHandler {
	return &LandHandlerImpl{uc, authUtil}
}

func (h *LandHandlerImpl) CreateLand(c *gin.Context) {
	userId, err := h.authUtil.GetAuthUserID(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	var req dto.LandCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	land, err := h.uc.CreateLand(c, userId, &req)
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
	land, err := h.uc.GetLandByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
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
	lands, err := h.uc.GetLandByUserID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, lands)
}

func (h *LandHandlerImpl) GetAllLands(c *gin.Context) {
	lands, err := h.uc.GetAllLands(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, lands)
}

func (h *LandHandlerImpl) UpdateLand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	userId, err := h.authUtil.GetAuthUserID(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	var req dto.LandUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedLand, err := h.uc.UpdateLand(c, userId, id, &req)
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
	err = h.uc.DeleteLand(c, id)
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
	restoredLand, err := h.uc.RestoreLand(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredLand)

}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type AuthHandlerImpl struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(usecase usecase.AuthUsecase) AuthHandler {
	return &AuthHandlerImpl{usecase}
}

func (h *AuthHandlerImpl) Login(c *gin.Context) {
	var req dto.AuthDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	auth, err := h.usecase.Login(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, auth)
}

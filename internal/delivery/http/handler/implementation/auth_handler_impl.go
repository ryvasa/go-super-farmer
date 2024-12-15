package handler_implementation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type AuthHandlerImpl struct {
	uc usecase_interface.AuthUsecase
}

func NewAuthHandler(uc usecase_interface.AuthUsecase) handler_interface.AuthHandler {
	return &AuthHandlerImpl{uc}
}

func (h *AuthHandlerImpl) Login(c *gin.Context) {
	var req dto.AuthDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	auth, err := h.uc.Login(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, auth)
}

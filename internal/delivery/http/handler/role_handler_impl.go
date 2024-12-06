package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type RoleHandlerImpl struct {
	uc usecase.RoleUsecase
}

func NewRoleHandler(uc usecase.RoleUsecase) RoleHandler {
	return &RoleHandlerImpl{uc: uc}
}

func (h *RoleHandlerImpl) CreateRole(c *gin.Context) {
	var req dto.RoleCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	role, err := h.uc.CreateRole(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, role)
}

func (h *RoleHandlerImpl) GetAllRoles(c *gin.Context) {
	roles, err := h.uc.GetAllRoles(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, roles)
}

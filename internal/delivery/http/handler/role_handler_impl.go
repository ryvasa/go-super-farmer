package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type RoleHandlerImpl struct {
	usecase usecase.RoleUsecase
}

func NewRoleHandler(uc usecase.RoleUsecase) RoleHandler {
	return &RoleHandlerImpl{usecase: uc}
}

func (h *RoleHandlerImpl) CreateRole(c *gin.Context) {
	var req dto.RoleCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	role, err := h.usecase.CreateRole(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, role)
}

func (h *RoleHandlerImpl) GetAllRoles(c *gin.Context) {
	roles, err := h.usecase.GetAllRoles(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, roles)
}

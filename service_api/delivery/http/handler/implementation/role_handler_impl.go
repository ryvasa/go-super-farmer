package handler_implementation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type RoleHandlerImpl struct {
	uc usecase_interface.RoleUsecase
}

func NewRoleHandler(uc usecase_interface.RoleUsecase) handler_interface.RoleHandler {
	return &RoleHandlerImpl{uc: uc}
}

func (h *RoleHandlerImpl) CreateRole(c *gin.Context) {
	var req dto.RoleCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	role, err := h.uc.CreateRole(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, role)
}

func (h *RoleHandlerImpl) GetAllRoles(c *gin.Context) {
	roles, err := h.uc.GetAllRoles(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, roles)
}

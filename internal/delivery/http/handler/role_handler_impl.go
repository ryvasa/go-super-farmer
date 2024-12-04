package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
)

type RoleHandlerImpl struct {
	uc usecase.RoleUsecase
}

func NewRoleHandler(uc usecase.RoleUsecase) RoleHandler {
	return &RoleHandlerImpl{uc: uc}
}

func (h *RoleHandlerImpl) CreateRole(c *gin.Context) {
	var role domain.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.CreateRole(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, role)
}

func (h *RoleHandlerImpl) GetAllRoles(c *gin.Context) {
	roles, err := h.uc.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

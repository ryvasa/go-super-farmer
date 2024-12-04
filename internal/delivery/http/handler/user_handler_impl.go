package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
)

type UserHandlerImpl struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) UserHandler {
	return &UserHandlerImpl{uc: uc}
}

func (h *UserHandlerImpl) RegisterUser(c *gin.Context) {
	var req dto.UserCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	createdUser, err := h.uc.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, createdUser)
}

func (h *UserHandlerImpl) GetOneUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	user, err := h.uc.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, user)
}

func (h *UserHandlerImpl) GetAllUsers(c *gin.Context) {
	users, err := h.uc.GetAllUsers()
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, users)
}

func (h *UserHandlerImpl) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	var req dto.UserUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	updatedUser, err := h.uc.UpdateUser(id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedUser)
}

func (h *UserHandlerImpl) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	if err := h.uc.DeleteUser(id); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandlerImpl) RestoreUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	restoredUser, err := h.uc.RestoreUser(id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredUser)
}

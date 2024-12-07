package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	createdUser, err := h.uc.Register(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, createdUser)
}

func (h *UserHandlerImpl) GetOneUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	user, err := h.uc.GetUserByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, user)
}

func (h *UserHandlerImpl) GetAllUsers(c *gin.Context) {
	users, err := h.uc.GetAllUsers(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, users)
}

func (h *UserHandlerImpl) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	var req dto.UserUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedUser, err := h.uc.UpdateUser(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedUser)
}

func (h *UserHandlerImpl) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	if err := h.uc.DeleteUser(c, id); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandlerImpl) RestoreUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	restoredUser, err := h.uc.RestoreUser(c, id)
	if err != nil {
		utils.ErrorResponse(c, utils.NewInternalError(err.Error()))
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredUser)
}

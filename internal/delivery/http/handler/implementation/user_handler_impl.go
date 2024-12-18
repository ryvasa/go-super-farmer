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

type UserHandlerImpl struct {
	uc     usecase_interface.UserUsecase
	ucAuth usecase_interface.AuthUsecase
}

func NewUserHandler(uc usecase_interface.UserUsecase, ucAuth usecase_interface.AuthUsecase) handler_interface.UserHandler {
	return &UserHandlerImpl{uc, ucAuth}
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
	err = h.ucAuth.SendOTP(c, &dto.AuthSendDTO{
		Email: createdUser.Email,
	})
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
	pagination, err := utils.GetPaginationParams(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	users, err := h.uc.GetAllUsers(c, pagination)
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
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredUser)
}

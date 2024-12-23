package handler_implementation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/service_api/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type ProvinceHandlerImpl struct {
	uc usecase_interface.ProvinceUsecase
}

func NewProvinceHandler(uc usecase_interface.ProvinceUsecase) handler_interface.ProvinceHandler {
	return &ProvinceHandlerImpl{uc}
}

func (h *ProvinceHandlerImpl) CreateProvince(c *gin.Context) {
	req := new(dto.ProvinceCreateDTO)
	if err := c.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	province, err := h.uc.CreateProvince(c, req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, province)
}

func (h *ProvinceHandlerImpl) GetAllProvinces(c *gin.Context) {
	provinces, err := h.uc.GetAllProvinces(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, provinces)
}

func (h *ProvinceHandlerImpl) GetProvinceByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	province, err := h.uc.GetProvinceByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, province)
}

func (h *ProvinceHandlerImpl) UpdateProvince(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	req := new(dto.ProvinceUpdateDTO)
	if err := c.ShouldBindJSON(req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	province, err := h.uc.UpdateProvince(c, id, req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, province)
}

func (h *ProvinceHandlerImpl) DeleteProvince(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	err = h.uc.DeleteProvince(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Province deleted successfully"})
}

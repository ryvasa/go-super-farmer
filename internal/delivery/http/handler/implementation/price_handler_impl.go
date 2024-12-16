package handler_implementation

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceHandlerImpl struct {
	uc usecase_interface.PriceUsecase
}

func NewPriceHandler(uc usecase_interface.PriceUsecase) handler_interface.PriceHandler {
	return &PriceHandlerImpl{uc}
}

func (h *PriceHandlerImpl) CreatePrice(c *gin.Context) {
	var req dto.PriceCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	createdPrice, err := h.uc.CreatePrice(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, createdPrice)
}

func (h *PriceHandlerImpl) GetAllPrices(c *gin.Context) {
	prices, err := h.uc.GetAllPrices(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, prices)
}

func (h *PriceHandlerImpl) GetPriceByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	user, err := h.uc.GetPriceByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, user)
}

func (h *PriceHandlerImpl) GetPricesByCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	prices, err := h.uc.GetPricesByCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, prices)
}

func (h *PriceHandlerImpl) GetPricesByRegionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	prices, err := h.uc.GetPricesByRegionID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, prices)
}

func (h *PriceHandlerImpl) UpdatePrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	req := dto.PriceUpdateDTO{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedPrice, err := h.uc.UpdatePrice(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedPrice)
}

func (h *PriceHandlerImpl) DeletePrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	if err := h.uc.DeletePrice(c, id); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Price deleted successfully"})
}

func (h *PriceHandlerImpl) RestorePrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	restoredPrice, err := h.uc.RestorePrice(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredPrice)
}

func (h *PriceHandlerImpl) GetPriceByCommodityIDAndRegionID(c *gin.Context) {
	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	regionID, err := uuid.Parse(c.Param("region_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	price, err := h.uc.GetPriceByCommodityIDAndRegionID(c, commodityID, regionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, price)
}

func (h *PriceHandlerImpl) GetPricesHistoryByCommodityIDAndRegionID(c *gin.Context) {
	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	regionID, err := uuid.Parse(c.Param("region_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	priceHistory, err := h.uc.GetPriceHistoryByCommodityIDAndRegionID(c, commodityID, regionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, priceHistory)
}

func (h *PriceHandlerImpl) DownloadPricesHistoryByCommodityIDAndRegionID(c *gin.Context) {
	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	regionID, err := uuid.Parse(c.Param("region_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	// Generate request ID
	requestID := fmt.Sprintf("%s_%s", commodityID, regionID)

	// Trigger Excel generation melalui RabbitMQ
	err = h.uc.DownloadPriceHistoryByCommodityIDAndRegionID(c, commodityID, regionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Tunggu beberapa detik untuk Excel generation
	time.Sleep(2 * time.Second)

	// Cek apakah file sudah tersedia
	filePath, exists := utils.GetFilePath(requestID)
	if !exists {
		utils.ErrorResponse(c, utils.NewInternalError("File not ready yet"))
		return
	}

	// Set header untuk download
	fileName := filepath.Base(filePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")

	// Kirim file
	c.File(filePath)

	// Cleanup
	defer func() {
		os.Remove(filePath)
		utils.RemoveFilePath(requestID)
	}()
}

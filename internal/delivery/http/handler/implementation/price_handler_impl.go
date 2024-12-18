package handler_implementation

import (
	"fmt"
	"net/http"
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

	priceParams := &dto.PriceParamsDTO{}

	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
			return
		}
		priceParams.StartDate = startDate
	}
	endDatestr := c.Query("end_date")
	if endDatestr != "" {
		endDate, err := time.Parse("2006-01-02", endDatestr)
		if err != nil {
			utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
			return
		}
		priceParams.EndDate = endDate
	}

	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	priceParams.CommodityID = commodityID
	regionID, err := uuid.Parse(c.Param("region_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	priceParams.RegionID = regionID

	err = h.uc.DownloadPriceHistoryByCommodityIDAndRegionID(c, priceParams)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Report generation in progress. Please check back in a few moments.",
		"download_url": fmt.Sprintf("http://localhost:8080/api/prices/history/commodity/%s/region/%s/download/file?start_date=%s&end_date=%s",
			commodityID, regionID, priceParams.StartDate.Format("2006-01-02"), priceParams.EndDate.Format("2006-01-02")),
	})

}

func (h *PriceHandlerImpl) GetPriceHistoryExcelFile(c *gin.Context) {
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

	startDateStr := c.Query("start_date")
	endDatestr := c.Query("end_date")

	// Get the latest excel file
	filePath := fmt.Sprintf("./public/reports/price_history_%s_%s_%s_%s_*.xlsx", commodityID, regionID, startDateStr, endDatestr)
	matches, err := filepath.Glob(filePath)
	if err != nil {
		utils.ErrorResponse(c, utils.NewInternalError("Error finding report file"))
		return
	}

	if len(matches) == 0 {
		utils.ErrorResponse(c, utils.NewNotFoundError("Report file not found"))
		return
	}

	// Get the latest file (assuming filename contains timestamp)
	latestFile := matches[len(matches)-1]

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(latestFile)))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Serve the file
	c.File(latestFile)
}

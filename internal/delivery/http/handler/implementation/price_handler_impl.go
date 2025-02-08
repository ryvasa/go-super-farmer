package handler_implementation

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	pb "github.com/ryvasa/go-super-farmer/proto/generated"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceHandlerImpl struct {
	uc           usecase_interface.PriceUsecase
	reportClient pb.ReportServiceClient
	minioClient  *minio.Client
}

func NewPriceHandler(uc usecase_interface.PriceUsecase, reportClient pb.ReportServiceClient, minioClient *minio.Client) handler_interface.PriceHandler {
	return &PriceHandlerImpl{uc, reportClient, minioClient}
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
	pagination, err := utils.GetPaginationParams(c)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	prices, err := h.uc.GetAllPrices(c, pagination)
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

func (h *PriceHandlerImpl) GetPricesByCityID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	prices, err := h.uc.GetPricesByCityID(c, id)
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

func (h *PriceHandlerImpl) GetPriceByCommodityIDAndCityID(c *gin.Context) {
	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	price, err := h.uc.GetPriceByCommodityIDAndCityID(c, commodityID, cityID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, price)
}

func (h *PriceHandlerImpl) GetPricesHistoryByCommodityIDAndCityID(c *gin.Context) {
	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	priceHistory, err := h.uc.GetPriceHistoryByCommodityIDAndCityID(c, commodityID, cityID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, priceHistory)
}

func (h *PriceHandlerImpl) GetReportPricesHistoryByCommodityIDAndCityID(c *gin.Context) {
	commodityID, err := uuid.Parse(c.Param("commodity_id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError("invalid commodity id"))
		return
	}

	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError("invalid city id"))
		return
	}

	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError("invalid start date"))
		return
	}

	endDate, err := time.Parse("2006-01-02", c.Query("end_date"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError("invalid end date"))
		return
	}

	// Prepare gRPC request
	req := &pb.PriceParams{
		CommodityId: commodityID.String(),
		CityId:      cityID,
		StartDate:   startDate.Format("2006-01-02"),
		EndDate:     endDate.Format("2006-01-02"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.reportClient.GetReportPrice(ctx, req)
	if err != nil {
		utils.ErrorResponse(c, utils.NewInternalError("failed to generate report"))
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/prices/history%s/download", resp.ReportUrl)
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"report_url": url,
	})
}

func (h *PriceHandlerImpl) DownloadFileReport(c *gin.Context) {
	fileName := c.Param("file_report")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File parameter is required"})
		return
	}

	bucket := c.Param("bucket")
	if bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bucket parameter is required"})
		return
	}

	// Check if file exists
	_, err := h.minioClient.StatObject(context.Background(), bucket, fileName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check file"})
		}
		return
	}

	// Download file from MinIO
	obj, err := h.minioClient.GetObject(context.Background(), bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}
	defer obj.Close()

	// Set header for downloading file
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	//Send file directly to client
	_, err = io.Copy(c.Writer, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error streaming file"})
		return
	}

	c.Writer.Flush()
}

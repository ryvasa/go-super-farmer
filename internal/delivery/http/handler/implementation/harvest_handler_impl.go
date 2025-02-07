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
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	pb "github.com/ryvasa/go-super-farmer/proto/generated"
	"github.com/ryvasa/go-super-farmer/utils"
)

type HarvestHandlerImpl struct {
	uc           usecase_interface.HarvestUsecase
	reportClient pb.ReportServiceClient
	minioClient  *minio.Client
}

func NewHarvestHandler(uc usecase_interface.HarvestUsecase,
	reportClient pb.ReportServiceClient, minioClient *minio.Client) handler_interface.HarvestHandler {
	return &HarvestHandlerImpl{uc, reportClient, minioClient}
}

func (h *HarvestHandlerImpl) CreateHarvest(c *gin.Context) {
	var req dto.HarvestCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvest, err := h.uc.CreateHarvest(c, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, harvest)
}

func (h *HarvestHandlerImpl) GetAllHarvest(c *gin.Context) {
	harvests, err := h.uc.GetAllHarvest(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvest, err := h.uc.GetHarvestByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvest)
}

func (h *HarvestHandlerImpl) GetHarvestByCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvests, err := h.uc.GetHarvestByCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByLandID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvests, err := h.uc.GetHarvestByLandID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByLandCommodityID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	harvests, err := h.uc.GetHarvestByLandCommodityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestByCityID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvests, err := h.uc.GetHarvestByCityID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) UpdateHarvest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	var req dto.HarvestUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	updatedHarvest, err := h.uc.UpdateHarvest(c, id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updatedHarvest)
}

func (h *HarvestHandlerImpl) DeleteHarvest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}

	err = h.uc.DeleteHarvest(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Harvest deleted successfully"})
}

func (h *HarvestHandlerImpl) RestoreHarvest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	restoredHarvest, err := h.uc.RestoreHarvest(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, restoredHarvest)
}

func (h *HarvestHandlerImpl) GetAllDeletedHarvest(c *gin.Context) {
	harvests, err := h.uc.GetAllDeletedHarvest(c)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvests)
}

func (h *HarvestHandlerImpl) GetHarvestDeletedByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
		return
	}
	harvest, err := h.uc.GetHarvestDeletedByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, harvest)
}
func (h *HarvestHandlerImpl) GetReportHarvestByLandCommodityID(c *gin.Context) {
	landCommodityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, utils.NewBadRequestError("invalid land commodity id"))
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
	req := &pb.HarvestParams{
		LandCommodityId: landCommodityID.String(),
		StartDate:       startDate.Format("2006-01-02"),
		EndDate:         endDate.Format("2006-01-02"),
	}

	fmt.Printf("Type of req: %T\n", req)
	// atau untuk field specific
	fmt.Printf("Type of LandCommodityId: %T\n", req.LandCommodityId)
	fmt.Printf("Type of StartDate: %T\n", req.StartDate)
	fmt.Printf("Type of EndDate: %T\n", req.EndDate)
	logrus.Log.Info(req)
	// Call report service

	// Call report service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.reportClient.GetReportHarvest(ctx, req)
	if err != nil {
		utils.ErrorResponse(c, utils.NewInternalError("failed to generate report"))
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/harvest%s/download", resp.ReportUrl)

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"report_url": url,
	})
}

func (h *HarvestHandlerImpl) DownloadFileReport(c *gin.Context) {
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

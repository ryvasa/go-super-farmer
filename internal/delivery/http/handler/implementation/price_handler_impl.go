package handler_implementation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_interface "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/interface"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	pb "github.com/ryvasa/go-super-farmer/proto/generated"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceHandlerImpl struct {
	uc           usecase_interface.PriceUsecase
	reportClient pb.ReportServiceClient
}

func NewPriceHandler(uc usecase_interface.PriceUsecase, reportClient pb.ReportServiceClient) handler_interface.PriceHandler {
	return &PriceHandlerImpl{uc, reportClient}
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

// TODO: change to gRPC
func (h *PriceHandlerImpl) DownloadPricesHistoryByCommodityIDAndCityID(c *gin.Context) {
	// logrus.Log.Info("Downloading price history")
	// priceParams := &dto.PriceParamsDTO{}

	// startDateStr := c.Query("start_date")
	// if startDateStr != "" {
	// 	startDate, err := time.Parse("2006-01-02", startDateStr)
	// 	if err != nil {
	// 		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
	// 		return
	// 	}
	// 	priceParams.StartDate = startDate
	// }
	// endDatestr := c.Query("end_date")
	// if endDatestr != "" {
	// 	endDate, err := time.Parse("2006-01-02", endDatestr)
	// 	if err != nil {
	// 		utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
	// 		return
	// 	}
	// 	priceParams.EndDate = endDate
	// }

	// commodityID, err := uuid.Parse(c.Param("commodity_id"))
	// if err != nil {
	// 	utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
	// 	return
	// }
	// priceParams.CommodityID = commodityID

	// cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	// if err != nil {
	// 	utils.ErrorResponse(c, utils.NewBadRequestError(err.Error()))
	// 	return
	// }
	// priceParams.CityID = cityID

	// logrus.Log.Info("Calling download price history usecase")
	// response, err := h.uc.DownloadPriceHistoryByCommodityIDAndCityID(c, priceParams)
	// if err != nil {
	// 	utils.ErrorResponse(c, err)
	// 	return
	// }
	// utils.SuccessResponse(c, http.StatusOK, response)
	// Parse request parameters
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
	fmt.Printf("Type of req: %T\n", req)
	// atau untuk field specific
	fmt.Printf("Type of CommodityId: %T\n", req.CommodityId)
	fmt.Printf("Type of CityId: %T\n", req.CityId)
	fmt.Printf("Type of StartDate: %T\n", req.StartDate)
	fmt.Printf("Type of EndDate: %T\n", req.EndDate)
	logrus.Log.Info(req)
	// Call report service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.reportClient.GetReportPrice(ctx, req)
	if err != nil {
		utils.ErrorResponse(c, utils.NewInternalError("failed to generate report"))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"report_url": resp.ReportUrl,
	})

}

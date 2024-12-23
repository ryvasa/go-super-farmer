package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/service_api/model/domain"
)

type ExcelInterface interface {
	CreatePriceHistoryReport(results []domain.PriceHistory, commodityName, regionName string, commodityID uuid.UUID, cityID int64, startDate, endDate time.Time) error
	CreateHarvestReport(results []domain.Harvest, commodityName, regionName, farmerName string, commodityID uuid.UUID, startDate, endDate time.Time) error
}

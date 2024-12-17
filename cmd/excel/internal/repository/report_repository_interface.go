package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
)

type ReportRepository interface {
	GetPriceHistoryReport(start, end time.Time, commodityID, regionID uuid.UUID) ([]domain.PriceHistory, error)
	GetHarvestReport(start, end time.Time, landCommodityID uuid.UUID) ([]domain.Harvest, error)
}

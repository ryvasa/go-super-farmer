package usecase

import "github.com/ryvasa/go-super-farmer/internal/model/domain"

type ExcelInterface interface {
	CreatePriceHistoryReport(results []domain.PriceHistory, commodityName, regionName string) error
	CreateHarvestReport(results []domain.Harvest, commodityName, regionName, farmerName string) error
}

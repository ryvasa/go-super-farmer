package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"gorm.io/gorm"
)

type ReportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepositoryImpl(db *gorm.DB) ReportRepository {
	return &ReportRepositoryImpl{db}
}

func (r *ReportRepositoryImpl) GetPriceHistoryReport(start, end time.Time, commodityID, regionID uuid.UUID) ([]domain.PriceHistory, error) {
	var results []domain.PriceHistory

	// Query current price
	var currentPrice domain.Price
	err := r.db.
		Preload("Commodity", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Description")
		}).
		Preload("Region", func(db *gorm.DB) *gorm.DB {
			return db.Omit("CreatedAt", "UpdatedAt", "DeletedAt")
		}).
		Preload("Region.Province").
		Preload("Region.City").
		Where("prices.commodity_id = ? AND prices.region_id = ?", commodityID, regionID).
		First(&currentPrice).Error
	if err != nil {
		return nil, err
	}

	// Konversi current price ke price history
	results = append(results, domain.PriceHistory{
		ID:          currentPrice.ID,
		CommodityID: currentPrice.CommodityID,
		RegionID:    currentPrice.RegionID,
		Price:       currentPrice.Price,
		Unit:        currentPrice.Unit,
		CreatedAt:   currentPrice.CreatedAt,
		Commodity:   currentPrice.Commodity,
		Region:      currentPrice.Region,
	})

	// Query price histories
	var histories []domain.PriceHistory
	err = r.db.Preload("Commodity").
		Preload("Region.City").
		Joins("JOIN commodities ON price_histories.commodity_id = commodities.id").
		Joins("JOIN regions ON price_histories.region_id = regions.id").
		Joins("JOIN cities ON regions.city_id = cities.id").
		Where("price_histories.commodity_id = ? AND price_histories.region_id = ? AND price_histories.deleted_at IS NULL",
			commodityID, regionID).
		Order("price_histories.created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}

	// Gabungkan hasil
	results = append(results, histories...)

	return results, nil
}

func (r *ReportRepositoryImpl) GetHarvestReport(start, end time.Time, landCommodityID uuid.UUID) ([]domain.Harvest, error) {
	var results []domain.Harvest
	err := r.db.
		Where("land_commodity_id = ? AND deleted_at IS NULL", landCommodityID).
		Preload("LandCommodity").
		Preload("Region").
		Preload("LandCommodity.Commodity").
		Preload("LandCommodity.Land").
		Preload("LandCommodity.Land.User").
		Preload("Region.Province").
		Preload("Region.City").
		Order("harvests.created_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

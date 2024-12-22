package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/service_api/model/dto"
	"gorm.io/gorm"
)

func GetPaginationParams(c *gin.Context) (*dto.PaginationDTO, error) {
	pagination := &dto.PaginationDTO{
		Limit: 10,                // default limit
		Page:  1,                 // default page
		Sort:  "created_at desc", // default sort
	}

	if err := c.ShouldBindQuery(pagination); err != nil {
		return nil, err
	}

	// Validate after binding
	if err := pagination.Validate(); err != nil {
		return nil, err
	}

	return pagination, nil
}

func ApplyFilters(filter *dto.PaginationFilterDTO) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.UserName != "" {
			db = db.Where("name LIKE ?", "%"+filter.UserName+"%")
		}
		if filter.CommodityName != "" {
			db = db.Where("name LIKE ?", "%"+filter.CommodityName+"%")
		}
		if filter.CityName != "" {
			db = db.Where("name LIKE ?", "%"+filter.CityName+"%")
		}
		if !filter.StartDate.IsZero() {
			db = db.Where("created_at >= ?", filter.StartDate)
		}
		if !filter.EndDate.IsZero() {
			db = db.Where("created_at <= ?", filter.EndDate)
		}
		return db
	}
}

func GetPaginationScope(pagination *dto.PaginationDTO) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := pagination.Page
		if page == 0 {
			page = 1
		}
		limit := pagination.Limit
		if limit == 0 {
			limit = 10
		}
		sort := pagination.Sort
		if sort == "" {
			sort = "created_at desc"
		}
		offset := (page - 1) * limit
		return db.
			Offset(offset).
			Limit(limit).
			Order(sort)
	}
}

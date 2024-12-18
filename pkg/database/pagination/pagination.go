package pagination

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type PaginationParams struct {
	Limit         int       `json:"limit,omitempty;query:limit"`
	Page          int       `json:"page,omitempty;query:page"`
	Sort          string    `json:"sort,omitempty;query:sort"`
	UserName      string    `json:"name,omitempty;query:user_name"`
	CommodityName string    `json:"commodity_name,omitempty;query:commodity_name"`
	StartDate     time.Time `json:"start_date,omitempty;query:start_date"`
	EndDate       time.Time `json:"end_date,omitempty;query:end_date"`
	TotalRows     int64     `json:"total_rows"`
	TotalPages    int       `json:"total_pages"`
	// Rows          interface{} `json:"rows"`
}

func (p *PaginationParams) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationParams) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *PaginationParams) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *PaginationParams) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at desc"
	}
	return p.Sort
}

func Paginate(value interface{}, pagination *PaginationParams, db *gorm.DB, field string) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		query := db.Offset(pagination.GetOffset()).
			Limit(pagination.GetLimit()).
			Order(pagination.GetSort())

		if pagination.UserName != "" {
			query = query.Where(field+" LIKE ?", "%"+pagination.UserName+"%")
		}
		if pagination.CommodityName != "" {
			query = query.Where(field+" LIKE ?", "%"+pagination.CommodityName+"%")
		}
		if !pagination.StartDate.IsZero() {
			query = query.Where("created_at >= ?", pagination.StartDate)
		}
		if !pagination.EndDate.IsZero() {
			query = query.Where("created_at <= ?", pagination.EndDate)
		}

		return query
	}
}

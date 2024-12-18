package dto

import "time"

type PaginationDTO struct {
	Limit         int       `form:"limit"`
	Page          int       `form:"page"`
	Sort          string    `form:"sort"`
	UserName      string    `form:"user_name"`
	CommodityName string    `form:"commodity_name"`
	StartDate     time.Time `form:"start_date"`
	EndDate       time.Time `form:"end_date"`
}

type PaginationResponseDTO struct {
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page"`
	Rows       interface{} `json:"rows"`
}

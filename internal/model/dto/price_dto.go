package dto

import (
	"github.com/google/uuid"
)

type PriceCreateDTO struct {
	CommodityID uuid.UUID `json:"commodity_id" validate:"required"`
	RegionID    uuid.UUID `json:"region_id" validate:"required"`
	Price       float64   `json:"price" validate:"required,min=1"`
}

type PriceUpdateDTO struct {
	CommodityID uuid.UUID `json:"commodity_id" validate:"required"`
	RegionID    uuid.UUID `json:"region_id" validate:"required"`
	Price       float64   `json:"price" validate:"required,min=1"`
}

type PriceResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	Price       float64   `json:"price"`
	CommodityID uuid.UUID `json:"-"`
	Commodity   struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"-"`
	} `json:"commodity"`
	RegionID uuid.UUID `json:"-"`
	Region   struct {
		ID         uuid.UUID `json:"id"`
		ProvinceID int64     `json:"-"`
		Province   struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"province"`
		CityID int64 `json:"-"`
		City   struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"city"`
	} `json:"region"`
}

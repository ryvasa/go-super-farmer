package dto

import "github.com/google/uuid"

type DemandCreateDTO struct {
	CommodityID uuid.UUID `json:"commodity_id" validate:"required"`
	RegionID    uuid.UUID `json:"region_id" validate:"required"`
	Quantity    float64   `json:"quantity" validate:"required,gte=0"`
}

type DemandUpdateDTO struct {
	CommodityID uuid.UUID `json:"commodity_id"`
	RegionID    uuid.UUID `json:"region_id"`
	Quantity    float64   `json:"quantity" validate:"gte=0"`
}

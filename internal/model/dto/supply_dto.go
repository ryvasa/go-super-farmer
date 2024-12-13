package dto

import "github.com/google/uuid"

type SupplyCreateDTO struct {
	CommodityID uuid.UUID `json:"commodity_id"`
	RegionID    uuid.UUID `json:"region_id"`
	Quantity    float64   `json:"quantity" validate:"required,gte=0"`
}

type SupplyUpdateDTO struct {
	CommodityID uuid.UUID `json:"commodity_id"`
	RegionID    uuid.UUID `json:"region_id"`
	Quantity    float64   `json:"quantity" validate:"required,gte=0"`
}

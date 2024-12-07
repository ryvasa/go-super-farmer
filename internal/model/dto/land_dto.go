package dto

import (
	"time"

	"github.com/google/uuid"
)

type LandCreateDTO struct {
	LandArea    int64  `json:"land_area" validate:"required,min=1,max=10000"`
	Certificate string `json:"certificate" validate:"required,min=1,max=255"`
}

type LandUpdateDTO struct {
	LandArea    int64  `json:"land_area" validate:"required,min=1,max=10000"`
	Certificate string `json:"certificate" validate:"required,min=1,max=255"`
}

type LandResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	LandArea    int64     `json:"land_area"`
	Certificate string    `json:"certificate"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

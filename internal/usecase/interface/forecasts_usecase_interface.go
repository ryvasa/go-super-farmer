package usecase_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
)

type ForecastsUsecase interface {
	GetForecastsByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) (*dto.ForecastsResponseDTO, error)
}

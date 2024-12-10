package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type PriceUsecaseImpl struct {
	priceRepo repository.PriceRepository
}

func NewPriceUsecase(priceRepo repository.PriceRepository) PriceUsecase {
	return &PriceUsecaseImpl{priceRepo}
}

func (u *PriceUsecaseImpl) Create(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error) {
	price := domain.Price{}
	// if err := utils.ValidateStruct(price); len(err) > 0 {
	// 	return nil, utils.NewValidationError(err)
	// }

	price.CommodityID = req.CommodityID
	price.RegionID = req.RegionID
	price.Price = req.Price
	price.ID = uuid.New()

	err := u.priceRepo.Create(ctx, &price)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdPrice, err := u.priceRepo.FindByID(ctx, price.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdPrice, nil
}

func (u *PriceUsecaseImpl) GetAll(ctx context.Context) (*[]dto.PriceResponseDTO, error) {
	return nil, nil
}

func (u *PriceUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	price, err := u.priceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return price, nil
}

func (u *PriceUsecaseImpl) GetByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]dto.PriceResponseDTO, error) {
	return nil, nil
}

func (u *PriceUsecaseImpl) GetByRegionID(ctx context.Context, regionID uuid.UUID) (*[]dto.PriceResponseDTO, error) {
	return nil, nil
}

func (u *PriceUsecaseImpl) Update(ctx context.Context, req *dto.PriceUpdateDTO) (*dto.PriceResponseDTO, error) {
	return nil, nil
}

func (u *PriceUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (u *PriceUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (*dto.PriceResponseDTO, error) {
	return nil, nil
}

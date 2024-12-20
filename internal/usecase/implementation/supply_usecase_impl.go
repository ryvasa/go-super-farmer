package usecase_implementation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/transaction"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/utils"
)

type SupplyUsecaseImpl struct {
	supplyRepo        repository_interface.SupplyRepository
	supplyHistoryRepo repository_interface.SupplyHistoryRepository
	commodityRepo     repository_interface.CommodityRepository
	cityRepo          repository_interface.CityRepository
	txManager         transaction.TransactionManager
}

func NewSupplyUsecase(supplyRepo repository_interface.SupplyRepository, supplyHistoryRepo repository_interface.SupplyHistoryRepository, commodityRepo repository_interface.CommodityRepository, cityRepo repository_interface.CityRepository, txManager transaction.TransactionManager) usecase_interface.SupplyUsecase {
	return &SupplyUsecaseImpl{
		supplyRepo,
		supplyHistoryRepo,
		commodityRepo,
		cityRepo,
		txManager,
	}
}
func (u *SupplyUsecaseImpl) CreateSupply(ctx context.Context, req *dto.SupplyCreateDTO) (*domain.Supply, error) {
	supply := domain.Supply{}
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	commodity, err := u.commodityRepo.FindByID(ctx, req.CommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	city, err := u.cityRepo.FindByID(ctx, req.CityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}
	supply.CommodityID = req.CommodityID
	supply.CityID = req.CityID
	supply.Quantity = req.Quantity
	supply.Commodity = commodity
	supply.City = city
	supply.ID = uuid.New()

	err = u.supplyRepo.Create(ctx, &supply)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdSupply, err := u.supplyRepo.FindByID(ctx, supply.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdSupply, nil
}

func (u *SupplyUsecaseImpl) GetAllSupply(ctx context.Context) ([]*domain.Supply, error) {
	supplies, err := u.supplyRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return supplies, nil
}

func (u *SupplyUsecaseImpl) GetSupplyByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error) {
	supply, err := u.supplyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError("internal error")
	}
	return supply, nil
}

func (u *SupplyUsecaseImpl) GetSupplyByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Supply, error) {
	supplies, err := u.supplyRepo.FindByCommodityID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return supplies, nil
}

func (u *SupplyUsecaseImpl) GetSupplyByCityID(ctx context.Context, cityID int64) ([]*domain.Supply, error) {
	supplies, err := u.supplyRepo.FindByCityID(ctx, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return supplies, nil
}

func (u *SupplyUsecaseImpl) UpdateSupply(ctx context.Context, id uuid.UUID, req *dto.SupplyUpdateDTO) (*domain.Supply, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	res := &domain.Supply{}
	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		logrus.Log.Info("starting supply update transaction")
		supply, err := u.supplyRepo.FindByID(txCtx, id)
		if err != nil {
			logrus.Log.Error(err, "failed to find supply")
			return utils.NewNotFoundError(err.Error())
		}
		historySupply := domain.SupplyHistory{
			ID:          uuid.New(),
			CommodityID: supply.CommodityID,
			Commodity:   supply.Commodity,
			CityID:      supply.CityID,
			City:        supply.City,
			Quantity:    supply.Quantity,
			Unit:        supply.Unit,
		}
		err = u.supplyHistoryRepo.Create(txCtx, &historySupply)
		if err != nil {
			logrus.Log.Error(err, "failed to create supply history")
			return err
		}
		logrus.Log.Info("supply history created")
		supply.Quantity = req.Quantity
		err = u.supplyRepo.Update(txCtx, id, supply)
		if err != nil {
			logrus.Log.Error(err, "failed to update supply")
			return err
		}
		updatedSupply, err := u.supplyRepo.FindByID(txCtx, id)
		if err != nil {
			logrus.Log.Error(err, "failed to find supply")
			return err
		}
		logrus.Log.Info("supply updated")
		res = updatedSupply
		return nil
	})
	logrus.Log.Info("supply update transaction completed")

	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return res, nil
}

func (u *SupplyUsecaseImpl) DeleteSupply(ctx context.Context, id uuid.UUID) error {
	_, err := u.supplyRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("supply not found")
	}
	err = u.supplyRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (u *SupplyUsecaseImpl) GetSupplyHistoryByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.SupplyHistory, error) {
	supplys, err := u.supplyHistoryRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	supply, err := u.supplyRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentSupply := &domain.SupplyHistory{
		ID:          supply.ID,
		CommodityID: supply.CommodityID,
		Commodity:   supply.Commodity,
		CityID:      supply.CityID,
		City:        supply.City,
		Quantity:    supply.Quantity,
		CreatedAt:   supply.CreatedAt,
		UpdatedAt:   supply.UpdatedAt,
		DeletedAt:   supply.DeletedAt,
	}

	allSupplyHistory := append(supplys, currentSupply)
	return allSupplyHistory, nil
}

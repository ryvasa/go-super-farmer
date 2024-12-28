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

type DemandUsecaseImpl struct {
	demandRepo        repository_interface.DemandRepository
	demandHistoryRepo repository_interface.DemandHistoryRepository
	commodityRepo     repository_interface.CommodityRepository
	cityRepo          repository_interface.CityRepository
	txManager         transaction.TransactionManager
}

func NewDemandUsecase(demandRepo repository_interface.DemandRepository, demandHistoryRepo repository_interface.DemandHistoryRepository, commodityRepo repository_interface.CommodityRepository, cityRepo repository_interface.CityRepository, txManager transaction.TransactionManager) usecase_interface.DemandUsecase {
	return &DemandUsecaseImpl{
		demandRepo,
		demandHistoryRepo,
		commodityRepo,
		cityRepo,
		txManager,
	}
}

func (u *DemandUsecaseImpl) CreateDemand(ctx context.Context, req *dto.DemandCreateDTO) (*domain.Demand, error) {
	demand := domain.Demand{}
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	_, err := u.commodityRepo.FindByID(ctx, req.CommodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}
	_, err = u.cityRepo.FindByID(ctx, req.CityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}
	demand.CommodityID = req.CommodityID
	demand.CityID = req.CityID
	demand.Quantity = req.Quantity
	demand.ID = uuid.New()

	err = u.demandRepo.Create(ctx, &demand)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdDemand, err := u.demandRepo.FindByID(ctx, demand.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdDemand, nil
}

func (u *DemandUsecaseImpl) GetAllDemands(ctx context.Context) ([]*domain.Demand, error) {
	demands, err := u.demandRepo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demands, nil
}
func (u *DemandUsecaseImpl) GetDemandByID(ctx context.Context, id uuid.UUID) (*domain.Demand, error) {
	demand, err := u.demandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demand, nil
}

func (u *DemandUsecaseImpl) GetDemandsByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Demand, error) {
	_, err := u.commodityRepo.FindByID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewNotFoundError("commodity not found")
	}

	demands, err := u.demandRepo.FindByCommodityID(ctx, commodityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demands, nil
}

func (u *DemandUsecaseImpl) GetDemandsByCityID(ctx context.Context, cityID int64) ([]*domain.Demand, error) {
	_, err := u.cityRepo.FindByID(ctx, cityID)
	if err != nil {
		return nil, utils.NewNotFoundError("city not found")
	}

	demands, err := u.demandRepo.FindByCityID(ctx, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return demands, nil
}

func (u *DemandUsecaseImpl) UpdateDemand(ctx context.Context, id uuid.UUID, req *dto.DemandUpdateDTO) (*domain.Demand, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, utils.NewValidationError(err)
	}
	res := &domain.Demand{}
	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		logrus.Log.Info("starting demand update transaction")
		demand, err := u.demandRepo.FindByID(txCtx, id)
		if err != nil {
			logrus.Log.Error(err, "failed to find demand")
			return utils.NewNotFoundError(err.Error())
		}
		historyDemand := domain.DemandHistory{
			ID:          uuid.New(),
			CommodityID: demand.CommodityID,
			Commodity:   demand.Commodity,
			CityID:      demand.CityID,
			City:        demand.City,
			Quantity:    demand.Quantity,
			Unit:        demand.Unit,
		}
		err = u.demandHistoryRepo.Create(txCtx, &historyDemand)
		if err != nil {
			logrus.Log.Error(err, "failed to create demand history")
			return err
		}
		logrus.Log.Info("demand history created")

		demand.Quantity = req.Quantity
		err = u.demandRepo.Update(txCtx, id, demand)
		if err != nil {
			logrus.Log.Error(err, "failed to update demand")
			return err
		}
		updatedDemand, err := u.demandRepo.FindByID(txCtx, id)
		if err != nil {
			logrus.Log.Error(err, "failed to find demand")
			return err
		}
		logrus.Log.Info("demand updated")
		res = updatedDemand
		return nil
	})
	logrus.Log.Info("demand update transaction completed")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return res, nil
}

func (u *DemandUsecaseImpl) DeleteDemand(ctx context.Context, id uuid.UUID) error {
	_, err := u.demandRepo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("demand not found")
	}
	err = u.demandRepo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}
	return nil
}

func (u *DemandUsecaseImpl) GetDemandHistoryByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.DemandHistory, error) {
	demands, err := u.demandHistoryRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	demand, err := u.demandRepo.FindByCommodityIDAndCityID(ctx, commodityID, cityID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	currentDemand := &domain.DemandHistory{
		ID:          demand.ID,
		CommodityID: demand.CommodityID,
		Commodity:   demand.Commodity,
		CityID:      demand.CityID,
		City:        demand.City,
		Quantity:    demand.Quantity,
		CreatedAt:   demand.CreatedAt,
		UpdatedAt:   demand.UpdatedAt,
		Unit:        demand.Unit,
		DeletedAt:   demand.DeletedAt,
	}

	allDemandHistory := append(demands, currentDemand)
	return allDemandHistory, nil
}

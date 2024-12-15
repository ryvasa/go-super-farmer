package usecase_implementation

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
)

type ProvinceUsecaseImpl struct {
	repo repository_interface.ProvinceRepository
}

func NewProvinceUsecase(repo repository_interface.ProvinceRepository) usecase_interface.ProvinceUsecase {
	return &ProvinceUsecaseImpl{repo}
}

func (uc *ProvinceUsecaseImpl) CreateProvince(ctx context.Context, req *dto.ProvinceCreateDTO) (*domain.Province, error) {
	province := domain.Province{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	province.Name = req.Name

	err := uc.repo.Create(ctx, &province)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	createdProvince, err := uc.repo.FindByID(ctx, province.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return createdProvince, nil
}

func (uc *ProvinceUsecaseImpl) GetAllProvinces(ctx context.Context) (*[]domain.Province, error) {
	provinces, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return provinces, nil
}
func (uc *ProvinceUsecaseImpl) GetProvinceByID(ctx context.Context, id int64) (*domain.Province, error) {
	province, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("province not found")
	}
	return province, nil
}

func (uc *ProvinceUsecaseImpl) UpdateProvince(ctx context.Context, id int64, req *dto.ProvinceUpdateDTO) (*domain.Province, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	province, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError("province not found")
	}

	province.Name = req.Name

	err = uc.repo.Update(ctx, id, province)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	updatedProvince, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return updatedProvince, nil
}

func (uc *ProvinceUsecaseImpl) DeleteProvince(ctx context.Context, id int64) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError("province not found")
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	return nil
}

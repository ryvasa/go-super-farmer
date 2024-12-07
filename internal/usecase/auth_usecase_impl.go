package usecase

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/utils"
)

type AuthUsecaseImpl struct {
	userRepo  repository.UserRepository
	tokenUtil utils.TokenUtil
}

func NewAuthUsecase(userRepo repository.UserRepository, tokenUtil utils.TokenUtil) AuthUsecase {
	return &AuthUsecaseImpl{userRepo, tokenUtil}
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, req *dto.AuthDTO) (*dto.AuthResponseDTO, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	if user == nil {
		return nil, utils.NewBadRequestError("invalid password or email")
	}
	res := utils.CheckPasswordHash(req.Password, user.Password)
	if res == false {
		return nil, utils.NewBadRequestError("invalid password or email")
	}

	token, err := u.tokenUtil.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return utils.AuthDtoFormat(user, token), nil
}

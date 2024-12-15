package usecase_implementation

import (
	"context"

	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/utils"
)

type AuthUsecaseImpl struct {
	userRepo repository_interface.UserRepository
	token    token.Token
	hash     utils.Hasher
}

func NewAuthUsecase(userRepo repository_interface.UserRepository, token token.Token, hash utils.Hasher) usecase_interface.AuthUsecase {
	return &AuthUsecaseImpl{userRepo, token, hash}
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, req *dto.AuthDTO) (*dto.AuthResponseDTO, error) {
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}

	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.NewBadRequestError("invalid password or email")
	}

	res := u.hash.ValidatePassword(req.Password, user.Password)
	if res == false {
		return nil, utils.NewBadRequestError("invalid password or email")
	}

	token, err := u.token.GenerateToken(user.ID, user.Role.Name)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	return utils.AuthDtoFormat(user, token), nil
}

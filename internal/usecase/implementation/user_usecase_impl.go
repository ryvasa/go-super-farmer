package usecase_implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	repository_interface "github.com/ryvasa/go-super-farmer/internal/repository/interface"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/utils"
)

type UserUsecaseImpl struct {
	repo  repository_interface.UserRepository
	hash  utils.Hasher
	cache cache.Cache
}

func NewUserUsecase(repo repository_interface.UserRepository, hash utils.Hasher, cache cache.Cache) usecase_interface.UserUsecase {
	return &UserUsecaseImpl{repo, hash, cache}
}

func (uc *UserUsecaseImpl) Register(ctx context.Context, req *dto.UserCreateDTO) (*dto.UserResponseDTO, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	user.Name = req.Name
	user.Email = req.Email
	user.ID = uuid.New()

	hashedPassword, err := uc.hash.HashPassword(req.Password)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	user.Password = hashedPassword
	err = uc.repo.Create(ctx, &user)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	createdUser, err := uc.repo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "user")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return utils.UserDtoFormat(createdUser), nil
}

func (uc *UserUsecaseImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	return utils.UserDtoFormat(user), nil
}

func (uc *UserUsecaseImpl) GetAllUsers(ctx context.Context, queryParams *dto.PaginationDTO) (*dto.PaginationResponseDTO, error) {
	if err := queryParams.Validate(); err != nil {
		return nil, utils.NewBadRequestError(err.Error())
	}

	cacheKey := fmt.Sprintf("user_list_page_%d_limit_%d_%s",
		queryParams.Page,
		queryParams.Limit,
		queryParams.Filter.UserName,
	)
	var response *dto.PaginationResponseDTO
	cached, err := uc.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		err := json.Unmarshal(cached, &response)
		if err != nil {
			logrus.Log.Errorf("Error: %v", err)
			return nil, utils.NewInternalError("invalid data")
		}

		// Convert data back to []*dto.UserResponseDTO
		if data, ok := response.Data.([]interface{}); ok {
			users := make([]*dto.UserResponseDTO, len(data))
			for i, item := range data {
				if userMap, ok := item.(map[string]interface{}); ok {
					userJSON, _ := json.Marshal(userMap)
					var user dto.UserResponseDTO
					json.Unmarshal(userJSON, &user)
					users[i] = &user
				}
			}
			response.Data = users
		}

		logrus.Log.Info("Cache hit")
		return response, nil
	}

	users, err := uc.repo.FindAll(ctx, queryParams)
	if err != nil {
		logrus.Log.Errorf("Error: %v", err)
		return nil, utils.NewInternalError(err.Error())
	}

	count, err := uc.repo.Count(ctx, &queryParams.Filter)
	if err != nil {
		logrus.Log.Errorf("Error: %v", err)
		return nil, utils.NewInternalError(err.Error())
	}

	usersDto := make([]*dto.UserResponseDTO, 0)
	for _, user := range users {
		usersDto = append(usersDto, utils.UserDtoFormat(user))
	}

	// Create response
	response = &dto.PaginationResponseDTO{
		TotalRows:  int64(count),
		TotalPages: int(math.Ceil(float64(count) / float64(queryParams.Limit))),
		Page:       queryParams.Page,
		Limit:      queryParams.Limit,
		Data:       usersDto,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	err = uc.cache.Set(ctx, cacheKey, responseJSON, 4*time.Minute)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return response, nil
}

func (uc *UserUsecaseImpl) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	user := domain.User{}
	if err := utils.ValidateStruct(req); len(err) > 0 {
		return nil, utils.NewValidationError(err)
	}
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Phone = &req.Phone

	if req.Password != "" {
		hashedPassword, err := uc.hash.HashPassword(req.Password)
		if err != nil {
			return nil, utils.NewInternalError(err.Error())
		}
		user.Password = hashedPassword
	}

	err = uc.repo.Update(ctx, id, &user)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	updatedUser, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "user")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return utils.UserDtoFormat(updatedUser), nil
}

func (uc *UserUsecaseImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return utils.NewNotFoundError(err.Error())
	}
	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "user")
	if err != nil {
		return utils.NewInternalError(err.Error())
	}

	return nil
}

func (uc *UserUsecaseImpl) RestoreUser(ctx context.Context, id uuid.UUID) (*dto.UserResponseDTO, error) {
	_, err := uc.repo.FindDeletedByID(ctx, id)
	if err != nil {
		return nil, utils.NewNotFoundError(err.Error())
	}
	err = uc.repo.Restore(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}
	restoredUser, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	err = uc.cache.DeleteByPattern(ctx, "user")
	if err != nil {
		return nil, utils.NewInternalError(err.Error())
	}

	return utils.UserDtoFormat(restoredUser), err
}

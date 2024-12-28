package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_pkg "github.com/ryvasa/go-super-farmer/pkg/mock"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/internal/model/dto"
	mock_repo "github.com/ryvasa/go-super-farmer/internal/repository/mock"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	usecase_interface "github.com/ryvasa/go-super-farmer/internal/usecase/interface"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

type CommodityRepoMock struct {
	Commodity *mock_repo.MockCommodityRepository
	Cache     *mock_pkg.MockCache
}

type CommodityIDs struct {
	CommodityID uuid.UUID
}

type CommodityMocks struct {
	Commodity        *domain.Commodity
	Commodities      []*domain.Commodity
	UpdatedCommodity *domain.Commodity
}

type CommodityDTOMock struct {
	Create     *dto.CommodityCreateDTO
	Update     *dto.CommodityUpdateDTO
	Pagination *dto.PaginationDTO
}

func CommodityUsecaseUtils(t *testing.T) (*CommodityIDs, *CommodityMocks, *CommodityDTOMock, *CommodityRepoMock, usecase_interface.CommodityUsecase, context.Context) {
	commodityID := uuid.New()

	ids := &CommodityIDs{
		CommodityID: commodityID,
	}

	mocks := &CommodityMocks{
		Commodity: &domain.Commodity{
			ID:          commodityID,
			Name:        "test commodity",
			Description: "test commodity description",
			Code:        "12345",
		},
		Commodities: []*domain.Commodity{
			{
				ID:          commodityID,
				Name:        "test commodity",
				Description: "test commodity description",
				Code:        "12345",
			},
		},
		UpdatedCommodity: &domain.Commodity{
			ID:          commodityID,
			Name:        "updated commodity",
			Description: "updated commodity description",
			Code:        "12345",
		},
	}

	updateName := "updated commodity"
	upadateDesc := "updated commodity description"
	code := "12345"

	dto := &CommodityDTOMock{
		Create: &dto.CommodityCreateDTO{
			Name:        "test commodity",
			Description: "test commodity description",
			Code:        "12345",
			Duration:    "20000",
		},
		Update: &dto.CommodityUpdateDTO{
			Name:        &updateName,
			Description: &upadateDesc,
			Code:        &code,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commodityRepo := mock_repo.NewMockCommodityRepository(ctrl)
	cache := mock_pkg.NewMockCache(ctrl)
	uc := usecase_implementation.NewCommodityUsecase(commodityRepo, cache)
	ctx := context.TODO()

	repo := &CommodityRepoMock{Commodity: commodityRepo, Cache: cache}

	return ids, mocks, dto, repo, uc, ctx
}

func TestCommodityUsecase_CreateCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should create commodity successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Commodity) error {
			c.ID = ids.CommodityID
			return nil
		}).Times(1)
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Cache.EXPECT().DeleteByPattern(ctx, "commodity").Return(nil).Times(1) // Tambahkan ekspektasi ini

		resp, err := uc.CreateCommodity(ctx, dtos.Create)

		assert.NoError(t, err)
		assert.Equal(t, dtos.Create.Name, resp.Name)
		assert.Equal(t, dtos.Create.Description, resp.Description)
		assert.Equal(t, dtos.Create.Code, resp.Code)
	})

	t.Run("should return error validation error", func(t *testing.T) {
		req := &dto.CommodityCreateDTO{}
		resp, err := uc.CreateCommodity(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when create commodity fails", func(t *testing.T) {
		repo.Commodity.EXPECT().Create(ctx, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		resp, err := uc.CreateCommodity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when get created commodity fails", func(t *testing.T) {
		repo.Commodity.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Commodity) error {
			c.ID = ids.CommodityID
			return nil
		}).Times(1)
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("find created commodity error")).Times(1)

		resp, err := uc.CreateCommodity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "find created commodity error")
	})

	t.Run("should return error when cache delete fails", func(t *testing.T) {
		repo.Commodity.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, c *domain.Commodity) error {
			c.ID = ids.CommodityID
			return nil
		}).Times(1)
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Cache.EXPECT().DeleteByPattern(ctx, "commodity").Return(utils.NewInternalError("cache delete error")).Times(1)

		resp, err := uc.CreateCommodity(ctx, dtos.Create)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "cache delete error")
	})
}

func TestCommodityUsecase_GetAllCommodities(t *testing.T) {
	_, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)

	queryParams := &dto.PaginationDTO{
		Limit: 10,
		Page:  1,
		Filter: dto.ParamFilterDTO{
			CommodityName: "test",
		},
	}

	cacheKey := fmt.Sprintf("commodity_list_page_%d_limit_%d_%s",
		queryParams.Page,
		queryParams.Limit,
		queryParams.Filter.CommodityName,
	)

	t.Run("should return commodities from cache", func(t *testing.T) {
		expectedResponse := &dto.PaginationResponseDTO{
			TotalRows:  10,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
			Data:       mocks.Commodities,
		}

		cachedJSON, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)

		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(cachedJSON, nil)

		resp, err := uc.GetAllCommodities(ctx, queryParams)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, resp)
	})

	t.Run("should return error when cache get fails", func(t *testing.T) {

		repo.Cache.EXPECT().Get(ctx, cacheKey).Return([]byte("invalid data"), nil).Times(1)

		resp, err := uc.GetAllCommodities(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid data")
	})

	t.Run("should return commodities from repository when cache miss", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil)
		repo.Commodity.EXPECT().FindAll(ctx, queryParams).Return(mocks.Commodities, nil)
		repo.Commodity.EXPECT().Count(ctx, &queryParams.Filter).Return(int64(1), nil)

		repo.Cache.EXPECT().Set(ctx, cacheKey, gomock.Any(), 4*time.Minute).Return(nil).Times(1)

		resp, err := uc.GetAllCommodities(ctx, queryParams)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.TotalRows)
		assert.Equal(t, mocks.Commodities, resp.Data)
	})

	t.Run("should return error when repository findAll fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil)
		repo.Commodity.EXPECT().FindAll(ctx, queryParams).Return(nil, fmt.Errorf("database error"))

		resp, err := uc.GetAllCommodities(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("should return error when repository count fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil)
		repo.Commodity.EXPECT().FindAll(ctx, queryParams).Return(mocks.Commodities, nil)
		repo.Commodity.EXPECT().Count(ctx, &queryParams.Filter).Return(int64(0), fmt.Errorf("count error"))

		resp, err := uc.GetAllCommodities(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "count error")
	})

	t.Run("should return error when pagination validation fails", func(t *testing.T) {
		invalidPagination := &dto.PaginationDTO{
			Page:  0,
			Limit: 0,
		}

		resp, err := uc.GetAllCommodities(ctx, invalidPagination)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "page must be greater than 0")
	})

	t.Run("should return error when set cache fails", func(t *testing.T) {
		repo.Cache.EXPECT().Get(ctx, cacheKey).Return(nil, nil)
		repo.Commodity.EXPECT().FindAll(ctx, queryParams).Return(mocks.Commodities, nil)
		repo.Commodity.EXPECT().Count(ctx, &queryParams.Filter).Return(int64(1), nil)

		repo.Cache.EXPECT().Set(ctx, cacheKey, gomock.Any(), 4*time.Minute).Return(utils.NewInternalError("cache set error")).Times(1)

		resp, err := uc.GetAllCommodities(ctx, queryParams)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "cache set error")
	})
}

func TestCommodityUsecase_GetCommodityById(t *testing.T) {
	ids, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)
	t.Run("should return commodity by id successfully", func(t *testing.T) {
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		resp, err := uc.GetCommodityById(ctx, ids.CommodityID)

		assert.NoError(t, err)
		assert.Equal(t, mocks.Commodity.Name, resp.Name)
		assert.Equal(t, mocks.Commodity.Description, resp.Description)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {

		repo.Commodity.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		resp, err := uc.GetCommodityById(ctx, ids.CommodityID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})
}
func TestCommodityUsecase_UpdateCommodity(t *testing.T) {
	ids, mocks, dtos, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should update commodity successfully", func(t *testing.T) {
		// Mock FindByID before update
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Mock Update commodity
		repo.Commodity.EXPECT().Update(ctx, ids.CommodityID, gomock.Any()).DoAndReturn(func(ctx context.Context, id uuid.UUID, c *domain.Commodity) error {
			c.ID = ids.CommodityID
			return nil
		}).Times(1)
		// Mock FindByID after update
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.UpdatedCommodity, nil).Times(1)
		// Mock Cache delete
		repo.Cache.EXPECT().DeleteByPattern(ctx, "commodity").Return(nil).Times(1)

		// Call the UpdateCommodity method
		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		// Assert success response
		assert.NotNil(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, *dtos.Update.Name, resp.Name)
		assert.Equal(t, *dtos.Update.Description, resp.Description)
		assert.Equal(t, *dtos.Update.Code, resp.Code)
	})

	t.Run("should return validation error", func(t *testing.T) {
		// Call the UpdateCommodity method with invalid data
		code := "1"
		req := &dto.CommodityUpdateDTO{
			Code: &code,
		}

		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, req)

		// Assert validation error
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "Validation failed")
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		// Mock FindByID failure
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		// Call the UpdateCommodity method
		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		// Assert error response
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when updating commodity", func(t *testing.T) {
		// Mock FindByID success
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Mock Update failure
		repo.Commodity.EXPECT().Update(ctx, ids.CommodityID, gomock.Any()).Return(utils.NewInternalError("internal error")).Times(1)

		// Call the UpdateCommodity method
		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		// Assert error response
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when finding updated commodity", func(t *testing.T) {
		// Mock FindByID before update
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Mock Update success
		repo.Commodity.EXPECT().Update(ctx, ids.CommodityID, gomock.Any()).Return(nil).Times(1)
		// Mock FindByID failure after update
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("find updated commodity error")).Times(1)

		// Call the UpdateCommodity method
		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		// Assert error response
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "find updated commodity error")
	})

	t.Run("should return error when deleting cache", func(t *testing.T) {
		// Mock FindByID before update
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Mock Update success
		repo.Commodity.EXPECT().Update(ctx, ids.CommodityID, gomock.Any()).Return(nil).Times(1)
		// Mock FindByID after update
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.UpdatedCommodity, nil).Times(1)
		// Mock Cache delete failure
		repo.Cache.EXPECT().DeleteByPattern(ctx, "commodity").Return(utils.NewInternalError("cache delete error")).Times(1)

		// Call the UpdateCommodity method
		resp, err := uc.UpdateCommodity(ctx, ids.CommodityID, dtos.Update)

		// Assert error response
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "cache delete error")
	})
}

func TestCommodityUsecase_DeleteCommodity(t *testing.T) {
	// Setup initial dependencies
	ids, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should delete commodity successfully", func(t *testing.T) {
		// Expect FindByID to return a valid commodity
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Expect Delete to execute without errors
		repo.Commodity.EXPECT().Delete(ctx, ids.CommodityID).Return(nil).Times(1)
		// Expect cache invalidation to succeed
		repo.Cache.EXPECT().DeleteByPattern(ctx, "commodity").Return(nil).Times(1)

		// Execute the DeleteCommodity method
		err := uc.DeleteCommodity(ctx, ids.CommodityID)

		// Assert no errors
		assert.NoError(t, err)
	})

	t.Run("should return error when commodity not found", func(t *testing.T) {
		// Simulate FindByID returning a "not found" error
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("commodity not found")).Times(1)

		// Execute the DeleteCommodity method
		err := uc.DeleteCommodity(ctx, ids.CommodityID)

		// Assert the error
		assert.Error(t, err)
		assert.EqualError(t, err, "commodity not found")
	})

	t.Run("should return error when deleting commodity fails", func(t *testing.T) {
		// Expect FindByID to succeed
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Simulate Delete returning an error
		repo.Commodity.EXPECT().Delete(ctx, ids.CommodityID).Return(utils.NewInternalError("internal error")).Times(1)

		// Execute the DeleteCommodity method
		err := uc.DeleteCommodity(ctx, ids.CommodityID)

		// Assert the error
		assert.Error(t, err)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when cache deletion fails", func(t *testing.T) {
		// Expect FindByID and Delete to succeed
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		repo.Commodity.EXPECT().Delete(ctx, ids.CommodityID).Return(nil).Times(1)
		// Simulate cache deletion failure
		repo.Cache.EXPECT().DeleteByPattern(ctx, "commodity").Return(fmt.Errorf("cache error")).Times(1)

		// Execute the DeleteCommodity method
		err := uc.DeleteCommodity(ctx, ids.CommodityID)

		// Assert the error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache error")
	})
}

func TestCommodityUsecase_RestoreCommodity(t *testing.T) {
	// Setup initial dependencies
	ids, mocks, _, repo, uc, ctx := CommodityUsecaseUtils(t)

	t.Run("should restore commodity successfully", func(t *testing.T) {
		// Expect FindDeletedByID to return the deleted commodity
		repo.Commodity.EXPECT().FindDeletedByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Expect Restore to succeed
		repo.Commodity.EXPECT().Restore(ctx, ids.CommodityID).Return(nil).Times(1)
		// Expect FindByID to return the restored commodity
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)

		// Execute the RestoreCommodity method
		resp, err := uc.RestoreCommodity(ctx, ids.CommodityID)

		// Assert no errors and validate the returned commodity
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, mocks.Commodity.Name, resp.Name)
		assert.Equal(t, mocks.Commodity.Description, resp.Description)
	})

	t.Run("should return error when commodity not found in deleted records", func(t *testing.T) {
		// Simulate FindDeletedByID returning a "not found" error
		repo.Commodity.EXPECT().FindDeletedByID(ctx, ids.CommodityID).Return(nil, utils.NewNotFoundError("deleted commodity not found")).Times(1)

		// Execute the RestoreCommodity method
		resp, err := uc.RestoreCommodity(ctx, ids.CommodityID)

		// Assert the error and validate the response
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "deleted commodity not found")
	})

	t.Run("should return error when restore operation fails", func(t *testing.T) {
		// Expect FindDeletedByID to return the deleted commodity
		repo.Commodity.EXPECT().FindDeletedByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Simulate Restore returning an internal error
		repo.Commodity.EXPECT().Restore(ctx, ids.CommodityID).Return(utils.NewInternalError("internal error")).Times(1)

		// Execute the RestoreCommodity method
		resp, err := uc.RestoreCommodity(ctx, ids.CommodityID)

		// Assert the error and validate the response
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "internal error")
	})

	t.Run("should return error when retrieving restored commodity fails", func(t *testing.T) {
		// Expect FindDeletedByID to return the deleted commodity
		repo.Commodity.EXPECT().FindDeletedByID(ctx, ids.CommodityID).Return(mocks.Commodity, nil).Times(1)
		// Expect Restore to succeed
		repo.Commodity.EXPECT().Restore(ctx, ids.CommodityID).Return(nil).Times(1)
		// Simulate FindByID returning an internal error
		repo.Commodity.EXPECT().FindByID(ctx, ids.CommodityID).Return(nil, utils.NewInternalError("retrieval error")).Times(1)

		// Execute the RestoreCommodity method
		resp, err := uc.RestoreCommodity(ctx, ids.CommodityID)

		// Assert the error and validate the response
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "retrieval error")
	})
}

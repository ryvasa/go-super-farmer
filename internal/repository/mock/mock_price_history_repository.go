// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/interface/price_history_repository_interface.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
)

// MockPriceHistoryRepository is a mock of PriceHistoryRepository interface.
type MockPriceHistoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPriceHistoryRepositoryMockRecorder
}

// MockPriceHistoryRepositoryMockRecorder is the mock recorder for MockPriceHistoryRepository.
type MockPriceHistoryRepositoryMockRecorder struct {
	mock *MockPriceHistoryRepository
}

// NewMockPriceHistoryRepository creates a new mock instance.
func NewMockPriceHistoryRepository(ctrl *gomock.Controller) *MockPriceHistoryRepository {
	mock := &MockPriceHistoryRepository{ctrl: ctrl}
	mock.recorder = &MockPriceHistoryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPriceHistoryRepository) EXPECT() *MockPriceHistoryRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPriceHistoryRepository) Create(ctx context.Context, priceHistory *domain.PriceHistory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, priceHistory)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPriceHistoryRepositoryMockRecorder) Create(ctx, priceHistory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPriceHistoryRepository)(nil).Create), ctx, priceHistory)
}

// FindByCommodityIDAndCityID mocks base method.
func (m *MockPriceHistoryRepository) FindByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.PriceHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCommodityIDAndCityID", ctx, commodityID, cityID)
	ret0, _ := ret[0].([]*domain.PriceHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCommodityIDAndCityID indicates an expected call of FindByCommodityIDAndCityID.
func (mr *MockPriceHistoryRepositoryMockRecorder) FindByCommodityIDAndCityID(ctx, commodityID, cityID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCommodityIDAndCityID", reflect.TypeOf((*MockPriceHistoryRepository)(nil).FindByCommodityIDAndCityID), ctx, commodityID, cityID)
}

// FindByID mocks base method.
func (m *MockPriceHistoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.PriceHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.PriceHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockPriceHistoryRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockPriceHistoryRepository)(nil).FindByID), ctx, id)
}
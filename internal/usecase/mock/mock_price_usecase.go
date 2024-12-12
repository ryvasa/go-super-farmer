// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/price_usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
	dto "github.com/ryvasa/go-super-farmer/internal/model/dto"
)

// MockPriceUsecase is a mock of PriceUsecase interface.
type MockPriceUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockPriceUsecaseMockRecorder
}

// MockPriceUsecaseMockRecorder is the mock recorder for MockPriceUsecase.
type MockPriceUsecaseMockRecorder struct {
	mock *MockPriceUsecase
}

// NewMockPriceUsecase creates a new mock instance.
func NewMockPriceUsecase(ctrl *gomock.Controller) *MockPriceUsecase {
	mock := &MockPriceUsecase{ctrl: ctrl}
	mock.recorder = &MockPriceUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPriceUsecase) EXPECT() *MockPriceUsecaseMockRecorder {
	return m.recorder
}

// CreatePrice mocks base method.
func (m *MockPriceUsecase) CreatePrice(ctx context.Context, req *dto.PriceCreateDTO) (*domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePrice", ctx, req)
	ret0, _ := ret[0].(*domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePrice indicates an expected call of CreatePrice.
func (mr *MockPriceUsecaseMockRecorder) CreatePrice(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePrice", reflect.TypeOf((*MockPriceUsecase)(nil).CreatePrice), ctx, req)
}

// DeletePrice mocks base method.
func (m *MockPriceUsecase) DeletePrice(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePrice", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePrice indicates an expected call of DeletePrice.
func (mr *MockPriceUsecaseMockRecorder) DeletePrice(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePrice", reflect.TypeOf((*MockPriceUsecase)(nil).DeletePrice), ctx, id)
}

// GetAllPrices mocks base method.
func (m *MockPriceUsecase) GetAllPrices(ctx context.Context) (*[]domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPrices", ctx)
	ret0, _ := ret[0].(*[]domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPrices indicates an expected call of GetAllPrices.
func (mr *MockPriceUsecaseMockRecorder) GetAllPrices(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPrices", reflect.TypeOf((*MockPriceUsecase)(nil).GetAllPrices), ctx)
}

// GetPriceByCommodityIDAndRegionID mocks base method.
func (m *MockPriceUsecase) GetPriceByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriceByCommodityIDAndRegionID", ctx, commodityID, regionID)
	ret0, _ := ret[0].(*domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPriceByCommodityIDAndRegionID indicates an expected call of GetPriceByCommodityIDAndRegionID.
func (mr *MockPriceUsecaseMockRecorder) GetPriceByCommodityIDAndRegionID(ctx, commodityID, regionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriceByCommodityIDAndRegionID", reflect.TypeOf((*MockPriceUsecase)(nil).GetPriceByCommodityIDAndRegionID), ctx, commodityID, regionID)
}

// GetPriceByID mocks base method.
func (m *MockPriceUsecase) GetPriceByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriceByID", ctx, id)
	ret0, _ := ret[0].(*domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPriceByID indicates an expected call of GetPriceByID.
func (mr *MockPriceUsecaseMockRecorder) GetPriceByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriceByID", reflect.TypeOf((*MockPriceUsecase)(nil).GetPriceByID), ctx, id)
}

// GetPriceHistoryByCommodityIDAndRegionID mocks base method.
func (m *MockPriceUsecase) GetPriceHistoryByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*[]domain.PriceHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriceHistoryByCommodityIDAndRegionID", ctx, commodityID, regionID)
	ret0, _ := ret[0].(*[]domain.PriceHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPriceHistoryByCommodityIDAndRegionID indicates an expected call of GetPriceHistoryByCommodityIDAndRegionID.
func (mr *MockPriceUsecaseMockRecorder) GetPriceHistoryByCommodityIDAndRegionID(ctx, commodityID, regionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriceHistoryByCommodityIDAndRegionID", reflect.TypeOf((*MockPriceUsecase)(nil).GetPriceHistoryByCommodityIDAndRegionID), ctx, commodityID, regionID)
}

// GetPricesByCommodityID mocks base method.
func (m *MockPriceUsecase) GetPricesByCommodityID(ctx context.Context, commodityID uuid.UUID) (*[]domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPricesByCommodityID", ctx, commodityID)
	ret0, _ := ret[0].(*[]domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPricesByCommodityID indicates an expected call of GetPricesByCommodityID.
func (mr *MockPriceUsecaseMockRecorder) GetPricesByCommodityID(ctx, commodityID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPricesByCommodityID", reflect.TypeOf((*MockPriceUsecase)(nil).GetPricesByCommodityID), ctx, commodityID)
}

// GetPricesByRegionID mocks base method.
func (m *MockPriceUsecase) GetPricesByRegionID(ctx context.Context, regionID uuid.UUID) (*[]domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPricesByRegionID", ctx, regionID)
	ret0, _ := ret[0].(*[]domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPricesByRegionID indicates an expected call of GetPricesByRegionID.
func (mr *MockPriceUsecaseMockRecorder) GetPricesByRegionID(ctx, regionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPricesByRegionID", reflect.TypeOf((*MockPriceUsecase)(nil).GetPricesByRegionID), ctx, regionID)
}

// RestorePrice mocks base method.
func (m *MockPriceUsecase) RestorePrice(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestorePrice", ctx, id)
	ret0, _ := ret[0].(*domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestorePrice indicates an expected call of RestorePrice.
func (mr *MockPriceUsecaseMockRecorder) RestorePrice(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestorePrice", reflect.TypeOf((*MockPriceUsecase)(nil).RestorePrice), ctx, id)
}

// UpdatePrice mocks base method.
func (m *MockPriceUsecase) UpdatePrice(ctx context.Context, id uuid.UUID, req *dto.PriceUpdateDTO) (*domain.Price, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePrice", ctx, id, req)
	ret0, _ := ret[0].(*domain.Price)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePrice indicates an expected call of UpdatePrice.
func (mr *MockPriceUsecaseMockRecorder) UpdatePrice(ctx, id, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePrice", reflect.TypeOf((*MockPriceUsecase)(nil).UpdatePrice), ctx, id, req)
}

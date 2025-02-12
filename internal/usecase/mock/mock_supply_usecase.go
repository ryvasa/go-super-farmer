// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/interface/supply_usecase_interface.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
	dto "github.com/ryvasa/go-super-farmer/internal/model/dto"
)

// MockSupplyUsecase is a mock of SupplyUsecase interface.
type MockSupplyUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockSupplyUsecaseMockRecorder
}

// MockSupplyUsecaseMockRecorder is the mock recorder for MockSupplyUsecase.
type MockSupplyUsecaseMockRecorder struct {
	mock *MockSupplyUsecase
}

// NewMockSupplyUsecase creates a new mock instance.
func NewMockSupplyUsecase(ctrl *gomock.Controller) *MockSupplyUsecase {
	mock := &MockSupplyUsecase{ctrl: ctrl}
	mock.recorder = &MockSupplyUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSupplyUsecase) EXPECT() *MockSupplyUsecaseMockRecorder {
	return m.recorder
}

// CreateSupply mocks base method.
func (m *MockSupplyUsecase) CreateSupply(ctx context.Context, req *dto.SupplyCreateDTO) (*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSupply", ctx, req)
	ret0, _ := ret[0].(*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSupply indicates an expected call of CreateSupply.
func (mr *MockSupplyUsecaseMockRecorder) CreateSupply(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSupply", reflect.TypeOf((*MockSupplyUsecase)(nil).CreateSupply), ctx, req)
}

// DeleteSupply mocks base method.
func (m *MockSupplyUsecase) DeleteSupply(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSupply", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSupply indicates an expected call of DeleteSupply.
func (mr *MockSupplyUsecaseMockRecorder) DeleteSupply(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSupply", reflect.TypeOf((*MockSupplyUsecase)(nil).DeleteSupply), ctx, id)
}

// GetAllSupply mocks base method.
func (m *MockSupplyUsecase) GetAllSupply(ctx context.Context) ([]*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSupply", ctx)
	ret0, _ := ret[0].([]*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSupply indicates an expected call of GetAllSupply.
func (mr *MockSupplyUsecaseMockRecorder) GetAllSupply(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSupply", reflect.TypeOf((*MockSupplyUsecase)(nil).GetAllSupply), ctx)
}

// GetSupplyByCityID mocks base method.
func (m *MockSupplyUsecase) GetSupplyByCityID(ctx context.Context, cityID int64) ([]*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupplyByCityID", ctx, cityID)
	ret0, _ := ret[0].([]*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupplyByCityID indicates an expected call of GetSupplyByCityID.
func (mr *MockSupplyUsecaseMockRecorder) GetSupplyByCityID(ctx, cityID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupplyByCityID", reflect.TypeOf((*MockSupplyUsecase)(nil).GetSupplyByCityID), ctx, cityID)
}

// GetSupplyByCommodityID mocks base method.
func (m *MockSupplyUsecase) GetSupplyByCommodityID(ctx context.Context, commodityID uuid.UUID) ([]*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupplyByCommodityID", ctx, commodityID)
	ret0, _ := ret[0].([]*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupplyByCommodityID indicates an expected call of GetSupplyByCommodityID.
func (mr *MockSupplyUsecaseMockRecorder) GetSupplyByCommodityID(ctx, commodityID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupplyByCommodityID", reflect.TypeOf((*MockSupplyUsecase)(nil).GetSupplyByCommodityID), ctx, commodityID)
}

// GetSupplyByID mocks base method.
func (m *MockSupplyUsecase) GetSupplyByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupplyByID", ctx, id)
	ret0, _ := ret[0].(*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupplyByID indicates an expected call of GetSupplyByID.
func (mr *MockSupplyUsecaseMockRecorder) GetSupplyByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupplyByID", reflect.TypeOf((*MockSupplyUsecase)(nil).GetSupplyByID), ctx, id)
}

// GetSupplyHistoryByCommodityIDAndCityID mocks base method.
func (m *MockSupplyUsecase) GetSupplyHistoryByCommodityIDAndCityID(ctx context.Context, commodityID uuid.UUID, cityID int64) ([]*domain.SupplyHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupplyHistoryByCommodityIDAndCityID", ctx, commodityID, cityID)
	ret0, _ := ret[0].([]*domain.SupplyHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupplyHistoryByCommodityIDAndCityID indicates an expected call of GetSupplyHistoryByCommodityIDAndCityID.
func (mr *MockSupplyUsecaseMockRecorder) GetSupplyHistoryByCommodityIDAndCityID(ctx, commodityID, cityID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupplyHistoryByCommodityIDAndCityID", reflect.TypeOf((*MockSupplyUsecase)(nil).GetSupplyHistoryByCommodityIDAndCityID), ctx, commodityID, cityID)
}

// UpdateSupply mocks base method.
func (m *MockSupplyUsecase) UpdateSupply(ctx context.Context, id uuid.UUID, req *dto.SupplyUpdateDTO) (*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSupply", ctx, id, req)
	ret0, _ := ret[0].(*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSupply indicates an expected call of UpdateSupply.
func (mr *MockSupplyUsecaseMockRecorder) UpdateSupply(ctx, id, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSupply", reflect.TypeOf((*MockSupplyUsecase)(nil).UpdateSupply), ctx, id, req)
}

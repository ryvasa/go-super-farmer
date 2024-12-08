// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/commodity_usecase.go

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

// MockCommodityUsecase is a mock of CommodityUsecase interface.
type MockCommodityUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockCommodityUsecaseMockRecorder
}

// MockCommodityUsecaseMockRecorder is the mock recorder for MockCommodityUsecase.
type MockCommodityUsecaseMockRecorder struct {
	mock *MockCommodityUsecase
}

// NewMockCommodityUsecase creates a new mock instance.
func NewMockCommodityUsecase(ctrl *gomock.Controller) *MockCommodityUsecase {
	mock := &MockCommodityUsecase{ctrl: ctrl}
	mock.recorder = &MockCommodityUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommodityUsecase) EXPECT() *MockCommodityUsecaseMockRecorder {
	return m.recorder
}

// CreateCommodity mocks base method.
func (m *MockCommodityUsecase) CreateCommodity(ctx context.Context, req *dto.CommodityCreateDTO) (*domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommodity", ctx, req)
	ret0, _ := ret[0].(*domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCommodity indicates an expected call of CreateCommodity.
func (mr *MockCommodityUsecaseMockRecorder) CreateCommodity(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommodity", reflect.TypeOf((*MockCommodityUsecase)(nil).CreateCommodity), ctx, req)
}

// DeleteCommodity mocks base method.
func (m *MockCommodityUsecase) DeleteCommodity(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCommodity", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommodity indicates an expected call of DeleteCommodity.
func (mr *MockCommodityUsecaseMockRecorder) DeleteCommodity(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommodity", reflect.TypeOf((*MockCommodityUsecase)(nil).DeleteCommodity), ctx, id)
}

// GetAllCommodities mocks base method.
func (m *MockCommodityUsecase) GetAllCommodities(ctx context.Context) (*[]domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCommodities", ctx)
	ret0, _ := ret[0].(*[]domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCommodities indicates an expected call of GetAllCommodities.
func (mr *MockCommodityUsecaseMockRecorder) GetAllCommodities(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCommodities", reflect.TypeOf((*MockCommodityUsecase)(nil).GetAllCommodities), ctx)
}

// GetCommodityById mocks base method.
func (m *MockCommodityUsecase) GetCommodityById(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommodityById", ctx, id)
	ret0, _ := ret[0].(*domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommodityById indicates an expected call of GetCommodityById.
func (mr *MockCommodityUsecaseMockRecorder) GetCommodityById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommodityById", reflect.TypeOf((*MockCommodityUsecase)(nil).GetCommodityById), ctx, id)
}

// RestoreCommodity mocks base method.
func (m *MockCommodityUsecase) RestoreCommodity(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreCommodity", ctx, id)
	ret0, _ := ret[0].(*domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestoreCommodity indicates an expected call of RestoreCommodity.
func (mr *MockCommodityUsecaseMockRecorder) RestoreCommodity(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreCommodity", reflect.TypeOf((*MockCommodityUsecase)(nil).RestoreCommodity), ctx, id)
}

// UpdateCommodity mocks base method.
func (m *MockCommodityUsecase) UpdateCommodity(ctx context.Context, id uuid.UUID, req *dto.CommodityUpdateDTO) (*domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCommodity", ctx, id, req)
	ret0, _ := ret[0].(*domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCommodity indicates an expected call of UpdateCommodity.
func (mr *MockCommodityUsecaseMockRecorder) UpdateCommodity(ctx, id, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCommodity", reflect.TypeOf((*MockCommodityUsecase)(nil).UpdateCommodity), ctx, id, req)
}

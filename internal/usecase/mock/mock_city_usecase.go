// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/interface/city_usecase_interface.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
	dto "github.com/ryvasa/go-super-farmer/internal/model/dto"
)

// MockCityUsecase is a mock of CityUsecase interface.
type MockCityUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockCityUsecaseMockRecorder
}

// MockCityUsecaseMockRecorder is the mock recorder for MockCityUsecase.
type MockCityUsecaseMockRecorder struct {
	mock *MockCityUsecase
}

// NewMockCityUsecase creates a new mock instance.
func NewMockCityUsecase(ctrl *gomock.Controller) *MockCityUsecase {
	mock := &MockCityUsecase{ctrl: ctrl}
	mock.recorder = &MockCityUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCityUsecase) EXPECT() *MockCityUsecaseMockRecorder {
	return m.recorder
}

// CreateCity mocks base method.
func (m *MockCityUsecase) CreateCity(ctx context.Context, req *dto.CityCreateDTO) (*domain.City, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCity", ctx, req)
	ret0, _ := ret[0].(*domain.City)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCity indicates an expected call of CreateCity.
func (mr *MockCityUsecaseMockRecorder) CreateCity(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCity", reflect.TypeOf((*MockCityUsecase)(nil).CreateCity), ctx, req)
}

// DeleteCity mocks base method.
func (m *MockCityUsecase) DeleteCity(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCity", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCity indicates an expected call of DeleteCity.
func (mr *MockCityUsecaseMockRecorder) DeleteCity(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCity", reflect.TypeOf((*MockCityUsecase)(nil).DeleteCity), ctx, id)
}

// GetAllCities mocks base method.
func (m *MockCityUsecase) GetAllCities(ctx context.Context) ([]*domain.City, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCities", ctx)
	ret0, _ := ret[0].([]*domain.City)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCities indicates an expected call of GetAllCities.
func (mr *MockCityUsecaseMockRecorder) GetAllCities(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCities", reflect.TypeOf((*MockCityUsecase)(nil).GetAllCities), ctx)
}

// GetCityByID mocks base method.
func (m *MockCityUsecase) GetCityByID(ctx context.Context, id int64) (*domain.City, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCityByID", ctx, id)
	ret0, _ := ret[0].(*domain.City)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCityByID indicates an expected call of GetCityByID.
func (mr *MockCityUsecaseMockRecorder) GetCityByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCityByID", reflect.TypeOf((*MockCityUsecase)(nil).GetCityByID), ctx, id)
}

// UpdateCity mocks base method.
func (m *MockCityUsecase) UpdateCity(ctx context.Context, id int64, req *dto.CityUpdateDTO) (*domain.City, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCity", ctx, id, req)
	ret0, _ := ret[0].(*domain.City)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCity indicates an expected call of UpdateCity.
func (mr *MockCityUsecaseMockRecorder) UpdateCity(ctx, id, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCity", reflect.TypeOf((*MockCityUsecase)(nil).UpdateCity), ctx, id, req)
}

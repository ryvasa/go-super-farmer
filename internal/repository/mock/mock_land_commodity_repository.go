// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/interface/land_commodity_repository_interface.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
	dto "github.com/ryvasa/go-super-farmer/internal/model/dto"
)

// MockLandCommodityRepository is a mock of LandCommodityRepository interface.
type MockLandCommodityRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLandCommodityRepositoryMockRecorder
}

// MockLandCommodityRepositoryMockRecorder is the mock recorder for MockLandCommodityRepository.
type MockLandCommodityRepositoryMockRecorder struct {
	mock *MockLandCommodityRepository
}

// NewMockLandCommodityRepository creates a new mock instance.
func NewMockLandCommodityRepository(ctrl *gomock.Controller) *MockLandCommodityRepository {
	mock := &MockLandCommodityRepository{ctrl: ctrl}
	mock.recorder = &MockLandCommodityRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLandCommodityRepository) EXPECT() *MockLandCommodityRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockLandCommodityRepository) Create(ctx context.Context, landCommodity *domain.LandCommodity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, landCommodity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockLandCommodityRepositoryMockRecorder) Create(ctx, landCommodity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLandCommodityRepository)(nil).Create), ctx, landCommodity)
}

// Delete mocks base method.
func (m *MockLandCommodityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockLandCommodityRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLandCommodityRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockLandCommodityRepository) FindAll(ctx context.Context) ([]*domain.LandCommodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*domain.LandCommodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockLandCommodityRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockLandCommodityRepository)(nil).FindAll), ctx)
}

// FindByCommodityID mocks base method.
func (m *MockLandCommodityRepository) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCommodityID", ctx, id)
	ret0, _ := ret[0].([]*domain.LandCommodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCommodityID indicates an expected call of FindByCommodityID.
func (mr *MockLandCommodityRepositoryMockRecorder) FindByCommodityID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCommodityID", reflect.TypeOf((*MockLandCommodityRepository)(nil).FindByCommodityID), ctx, id)
}

// FindByID mocks base method.
func (m *MockLandCommodityRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.LandCommodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockLandCommodityRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockLandCommodityRepository)(nil).FindByID), ctx, id)
}

// FindByLandID mocks base method.
func (m *MockLandCommodityRepository) FindByLandID(ctx context.Context, id uuid.UUID) ([]*domain.LandCommodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByLandID", ctx, id)
	ret0, _ := ret[0].([]*domain.LandCommodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByLandID indicates an expected call of FindByLandID.
func (mr *MockLandCommodityRepositoryMockRecorder) FindByLandID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByLandID", reflect.TypeOf((*MockLandCommodityRepository)(nil).FindByLandID), ctx, id)
}

// FindDeletedByID mocks base method.
func (m *MockLandCommodityRepository) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.LandCommodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeletedByID", ctx, id)
	ret0, _ := ret[0].(*domain.LandCommodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeletedByID indicates an expected call of FindDeletedByID.
func (mr *MockLandCommodityRepositoryMockRecorder) FindDeletedByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeletedByID", reflect.TypeOf((*MockLandCommodityRepository)(nil).FindDeletedByID), ctx, id)
}

// Restore mocks base method.
func (m *MockLandCommodityRepository) Restore(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockLandCommodityRepositoryMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockLandCommodityRepository)(nil).Restore), ctx, id)
}

// SumAllLandCommodityArea mocks base method.
func (m *MockLandCommodityRepository) SumAllLandCommodityArea(ctx context.Context, params *dto.LandAreaParamsDTO) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SumAllLandCommodityArea", ctx, params)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SumAllLandCommodityArea indicates an expected call of SumAllLandCommodityArea.
func (mr *MockLandCommodityRepositoryMockRecorder) SumAllLandCommodityArea(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SumAllLandCommodityArea", reflect.TypeOf((*MockLandCommodityRepository)(nil).SumAllLandCommodityArea), ctx, params)
}

// SumLandAreaByCommodityID mocks base method.
func (m *MockLandCommodityRepository) SumLandAreaByCommodityID(ctx context.Context, id uuid.UUID) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SumLandAreaByCommodityID", ctx, id)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SumLandAreaByCommodityID indicates an expected call of SumLandAreaByCommodityID.
func (mr *MockLandCommodityRepositoryMockRecorder) SumLandAreaByCommodityID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SumLandAreaByCommodityID", reflect.TypeOf((*MockLandCommodityRepository)(nil).SumLandAreaByCommodityID), ctx, id)
}

// SumLandAreaByLandID mocks base method.
func (m *MockLandCommodityRepository) SumLandAreaByLandID(ctx context.Context, id uuid.UUID) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SumLandAreaByLandID", ctx, id)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SumLandAreaByLandID indicates an expected call of SumLandAreaByLandID.
func (mr *MockLandCommodityRepositoryMockRecorder) SumLandAreaByLandID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SumLandAreaByLandID", reflect.TypeOf((*MockLandCommodityRepository)(nil).SumLandAreaByLandID), ctx, id)
}

// Update mocks base method.
func (m *MockLandCommodityRepository) Update(ctx context.Context, id uuid.UUID, landCommodity *domain.LandCommodity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, landCommodity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockLandCommodityRepositoryMockRecorder) Update(ctx, id, landCommodity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLandCommodityRepository)(nil).Update), ctx, id, landCommodity)
}
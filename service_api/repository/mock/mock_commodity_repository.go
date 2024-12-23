// Code generated by MockGen. DO NOT EDIT.
// Source: service_api/repository/interface/commodity_repository_interface.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/service_api/model/domain"
	dto "github.com/ryvasa/go-super-farmer/service_api/model/dto"
)

// MockCommodityRepository is a mock of CommodityRepository interface.
type MockCommodityRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCommodityRepositoryMockRecorder
}

// MockCommodityRepositoryMockRecorder is the mock recorder for MockCommodityRepository.
type MockCommodityRepositoryMockRecorder struct {
	mock *MockCommodityRepository
}

// NewMockCommodityRepository creates a new mock instance.
func NewMockCommodityRepository(ctrl *gomock.Controller) *MockCommodityRepository {
	mock := &MockCommodityRepository{ctrl: ctrl}
	mock.recorder = &MockCommodityRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommodityRepository) EXPECT() *MockCommodityRepositoryMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockCommodityRepository) Count(ctx context.Context, filter *dto.PaginationFilterDTO) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, filter)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockCommodityRepositoryMockRecorder) Count(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockCommodityRepository)(nil).Count), ctx, filter)
}

// Create mocks base method.
func (m *MockCommodityRepository) Create(ctx context.Context, land *domain.Commodity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, land)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCommodityRepositoryMockRecorder) Create(ctx, land interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCommodityRepository)(nil).Create), ctx, land)
}

// Delete mocks base method.
func (m *MockCommodityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCommodityRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCommodityRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockCommodityRepository) FindAll(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, params)
	ret0, _ := ret[0].([]*domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockCommodityRepositoryMockRecorder) FindAll(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockCommodityRepository)(nil).FindAll), ctx, params)
}

// FindByID mocks base method.
func (m *MockCommodityRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockCommodityRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockCommodityRepository)(nil).FindByID), ctx, id)
}

// FindDeletedByID mocks base method.
func (m *MockCommodityRepository) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Commodity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeletedByID", ctx, id)
	ret0, _ := ret[0].(*domain.Commodity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeletedByID indicates an expected call of FindDeletedByID.
func (mr *MockCommodityRepositoryMockRecorder) FindDeletedByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeletedByID", reflect.TypeOf((*MockCommodityRepository)(nil).FindDeletedByID), ctx, id)
}

// Restore mocks base method.
func (m *MockCommodityRepository) Restore(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockCommodityRepositoryMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockCommodityRepository)(nil).Restore), ctx, id)
}

// Update mocks base method.
func (m *MockCommodityRepository) Update(ctx context.Context, id uuid.UUID, land *domain.Commodity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, land)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCommodityRepositoryMockRecorder) Update(ctx, id, land interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCommodityRepository)(nil).Update), ctx, id, land)
}
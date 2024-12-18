// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/region_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
)

// MockRegionRepository is a mock of RegionRepository interface.
type MockRegionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRegionRepositoryMockRecorder
}

// MockRegionRepositoryMockRecorder is the mock recorder for MockRegionRepository.
type MockRegionRepositoryMockRecorder struct {
	mock *MockRegionRepository
}

// NewMockRegionRepository creates a new mock instance.
func NewMockRegionRepository(ctrl *gomock.Controller) *MockRegionRepository {
	mock := &MockRegionRepository{ctrl: ctrl}
	mock.recorder = &MockRegionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegionRepository) EXPECT() *MockRegionRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRegionRepository) Create(ctx context.Context, region *domain.Region) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRegionRepositoryMockRecorder) Create(ctx, region interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRegionRepository)(nil).Create), ctx, region)
}

// Delete mocks base method.
func (m *MockRegionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRegionRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRegionRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockRegionRepository) FindAll(ctx context.Context) ([]*domain.Region, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*domain.Region)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockRegionRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockRegionRepository)(nil).FindAll), ctx)
}

// FindByID mocks base method.
func (m *MockRegionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Region, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Region)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockRegionRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockRegionRepository)(nil).FindByID), ctx, id)
}

// FindByProvinceID mocks base method.
func (m *MockRegionRepository) FindByProvinceID(ctx context.Context, id int64) ([]*domain.Region, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByProvinceID", ctx, id)
	ret0, _ := ret[0].([]*domain.Region)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByProvinceID indicates an expected call of FindByProvinceID.
func (mr *MockRegionRepositoryMockRecorder) FindByProvinceID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByProvinceID", reflect.TypeOf((*MockRegionRepository)(nil).FindByProvinceID), ctx, id)
}

// FindDeleted mocks base method.
func (m *MockRegionRepository) FindDeleted(ctx context.Context, id uuid.UUID) (*domain.Region, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeleted", ctx, id)
	ret0, _ := ret[0].(*domain.Region)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeleted indicates an expected call of FindDeleted.
func (mr *MockRegionRepositoryMockRecorder) FindDeleted(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeleted", reflect.TypeOf((*MockRegionRepository)(nil).FindDeleted), ctx, id)
}

// Restore mocks base method.
func (m *MockRegionRepository) Restore(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockRegionRepositoryMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockRegionRepository)(nil).Restore), ctx, id)
}

// Update mocks base method.
func (m *MockRegionRepository) Update(ctx context.Context, id uuid.UUID, region *domain.Region) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, region)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRegionRepositoryMockRecorder) Update(ctx, id, region interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRegionRepository)(nil).Update), ctx, id, region)
}

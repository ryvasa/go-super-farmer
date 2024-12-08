// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/land_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
)

// MockLandRepository is a mock of LandRepository interface.
type MockLandRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLandRepositoryMockRecorder
}

// MockLandRepositoryMockRecorder is the mock recorder for MockLandRepository.
type MockLandRepositoryMockRecorder struct {
	mock *MockLandRepository
}

// NewMockLandRepository creates a new mock instance.
func NewMockLandRepository(ctrl *gomock.Controller) *MockLandRepository {
	mock := &MockLandRepository{ctrl: ctrl}
	mock.recorder = &MockLandRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLandRepository) EXPECT() *MockLandRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockLandRepository) Create(ctx context.Context, land *domain.Land) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, land)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockLandRepositoryMockRecorder) Create(ctx, land interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLandRepository)(nil).Create), ctx, land)
}

// Delete mocks base method.
func (m *MockLandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockLandRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLandRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockLandRepository) FindAll(ctx context.Context) (*[]domain.Land, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].(*[]domain.Land)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockLandRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockLandRepository)(nil).FindAll), ctx)
}

// FindByID mocks base method.
func (m *MockLandRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Land)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockLandRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockLandRepository)(nil).FindByID), ctx, id)
}

// FindByUserID mocks base method.
func (m *MockLandRepository) FindByUserID(ctx context.Context, id uuid.UUID) (*[]domain.Land, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", ctx, id)
	ret0, _ := ret[0].(*[]domain.Land)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID.
func (mr *MockLandRepositoryMockRecorder) FindByUserID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockLandRepository)(nil).FindByUserID), ctx, id)
}

// FindDeletedByID mocks base method.
func (m *MockLandRepository) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeletedByID", ctx, id)
	ret0, _ := ret[0].(*domain.Land)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeletedByID indicates an expected call of FindDeletedByID.
func (mr *MockLandRepositoryMockRecorder) FindDeletedByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeletedByID", reflect.TypeOf((*MockLandRepository)(nil).FindDeletedByID), ctx, id)
}

// Restore mocks base method.
func (m *MockLandRepository) Restore(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockLandRepositoryMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockLandRepository)(nil).Restore), ctx, id)
}

// Update mocks base method.
func (m *MockLandRepository) Update(ctx context.Context, id uuid.UUID, land *domain.Land) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, land)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockLandRepositoryMockRecorder) Update(ctx, id, land interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLandRepository)(nil).Update), ctx, id, land)
}
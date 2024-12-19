// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/interface/province_repository_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
)

// MockProvinceRepository is a mock of ProvinceRepository interface.
type MockProvinceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProvinceRepositoryMockRecorder
}

// MockProvinceRepositoryMockRecorder is the mock recorder for MockProvinceRepository.
type MockProvinceRepositoryMockRecorder struct {
	mock *MockProvinceRepository
}

// NewMockProvinceRepository creates a new mock instance.
func NewMockProvinceRepository(ctrl *gomock.Controller) *MockProvinceRepository {
	mock := &MockProvinceRepository{ctrl: ctrl}
	mock.recorder = &MockProvinceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvinceRepository) EXPECT() *MockProvinceRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProvinceRepository) Create(ctx context.Context, province *domain.Province) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, province)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockProvinceRepositoryMockRecorder) Create(ctx, province interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProvinceRepository)(nil).Create), ctx, province)
}

// Delete mocks base method.
func (m *MockProvinceRepository) Delete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProvinceRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProvinceRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockProvinceRepository) FindAll(ctx context.Context) ([]*domain.Province, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*domain.Province)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockProvinceRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockProvinceRepository)(nil).FindAll), ctx)
}

// FindByID mocks base method.
func (m *MockProvinceRepository) FindByID(ctx context.Context, id int64) (*domain.Province, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Province)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockProvinceRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockProvinceRepository)(nil).FindByID), ctx, id)
}

// Update mocks base method.
func (m *MockProvinceRepository) Update(ctx context.Context, id int64, province *domain.Province) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, province)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockProvinceRepositoryMockRecorder) Update(ctx, id, province interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProvinceRepository)(nil).Update), ctx, id, province)
}

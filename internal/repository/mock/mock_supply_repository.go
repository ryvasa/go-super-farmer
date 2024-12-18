// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/supply_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
)

// MockSupplyRepository is a mock of SupplyRepository interface.
type MockSupplyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSupplyRepositoryMockRecorder
}

// MockSupplyRepositoryMockRecorder is the mock recorder for MockSupplyRepository.
type MockSupplyRepositoryMockRecorder struct {
	mock *MockSupplyRepository
}

// NewMockSupplyRepository creates a new mock instance.
func NewMockSupplyRepository(ctrl *gomock.Controller) *MockSupplyRepository {
	mock := &MockSupplyRepository{ctrl: ctrl}
	mock.recorder = &MockSupplyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSupplyRepository) EXPECT() *MockSupplyRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSupplyRepository) Create(ctx context.Context, supply *domain.Supply) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, supply)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSupplyRepositoryMockRecorder) Create(ctx, supply interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSupplyRepository)(nil).Create), ctx, supply)
}

// Delete mocks base method.
func (m *MockSupplyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSupplyRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSupplyRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockSupplyRepository) FindAll(ctx context.Context) ([]*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockSupplyRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockSupplyRepository)(nil).FindAll), ctx)
}

// FindByCommodityID mocks base method.
func (m *MockSupplyRepository) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCommodityID", ctx, id)
	ret0, _ := ret[0].([]*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCommodityID indicates an expected call of FindByCommodityID.
func (mr *MockSupplyRepositoryMockRecorder) FindByCommodityID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCommodityID", reflect.TypeOf((*MockSupplyRepository)(nil).FindByCommodityID), ctx, id)
}

// FindByCommodityIDAndRegionID mocks base method.
func (m *MockSupplyRepository) FindByCommodityIDAndRegionID(ctx context.Context, commodityID, regionID uuid.UUID) (*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCommodityIDAndRegionID", ctx, commodityID, regionID)
	ret0, _ := ret[0].(*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCommodityIDAndRegionID indicates an expected call of FindByCommodityIDAndRegionID.
func (mr *MockSupplyRepositoryMockRecorder) FindByCommodityIDAndRegionID(ctx, commodityID, regionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCommodityIDAndRegionID", reflect.TypeOf((*MockSupplyRepository)(nil).FindByCommodityIDAndRegionID), ctx, commodityID, regionID)
}

// FindByID mocks base method.
func (m *MockSupplyRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockSupplyRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockSupplyRepository)(nil).FindByID), ctx, id)
}

// FindByRegionID mocks base method.
func (m *MockSupplyRepository) FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.Supply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByRegionID", ctx, id)
	ret0, _ := ret[0].([]*domain.Supply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByRegionID indicates an expected call of FindByRegionID.
func (mr *MockSupplyRepositoryMockRecorder) FindByRegionID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByRegionID", reflect.TypeOf((*MockSupplyRepository)(nil).FindByRegionID), ctx, id)
}

// Update mocks base method.
func (m *MockSupplyRepository) Update(ctx context.Context, id uuid.UUID, supply *domain.Supply) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, supply)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSupplyRepositoryMockRecorder) Update(ctx, id, supply interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSupplyRepository)(nil).Update), ctx, id, supply)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/interface/harvest_repository_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	domain "github.com/ryvasa/go-super-farmer/internal/model/domain"
)

// MockHarvestRepository is a mock of HarvestRepository interface.
type MockHarvestRepository struct {
	ctrl     *gomock.Controller
	recorder *MockHarvestRepositoryMockRecorder
}

// MockHarvestRepositoryMockRecorder is the mock recorder for MockHarvestRepository.
type MockHarvestRepositoryMockRecorder struct {
	mock *MockHarvestRepository
}

// NewMockHarvestRepository creates a new mock instance.
func NewMockHarvestRepository(ctrl *gomock.Controller) *MockHarvestRepository {
	mock := &MockHarvestRepository{ctrl: ctrl}
	mock.recorder = &MockHarvestRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHarvestRepository) EXPECT() *MockHarvestRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockHarvestRepository) Create(ctx context.Context, harvest *domain.Harvest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, harvest)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockHarvestRepositoryMockRecorder) Create(ctx, harvest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockHarvestRepository)(nil).Create), ctx, harvest)
}

// Delete mocks base method.
func (m *MockHarvestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockHarvestRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockHarvestRepository)(nil).Delete), ctx, id)
}

// FindAll mocks base method.
func (m *MockHarvestRepository) FindAll(ctx context.Context) ([]*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockHarvestRepositoryMockRecorder) FindAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockHarvestRepository)(nil).FindAll), ctx)
}

// FindAllDeleted mocks base method.
func (m *MockHarvestRepository) FindAllDeleted(ctx context.Context) ([]*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllDeleted", ctx)
	ret0, _ := ret[0].([]*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllDeleted indicates an expected call of FindAllDeleted.
func (mr *MockHarvestRepositoryMockRecorder) FindAllDeleted(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllDeleted", reflect.TypeOf((*MockHarvestRepository)(nil).FindAllDeleted), ctx)
}

// FindByCommodityID mocks base method.
func (m *MockHarvestRepository) FindByCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCommodityID", ctx, id)
	ret0, _ := ret[0].([]*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCommodityID indicates an expected call of FindByCommodityID.
func (mr *MockHarvestRepositoryMockRecorder) FindByCommodityID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCommodityID", reflect.TypeOf((*MockHarvestRepository)(nil).FindByCommodityID), ctx, id)
}

// FindByID mocks base method.
func (m *MockHarvestRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockHarvestRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockHarvestRepository)(nil).FindByID), ctx, id)
}

// FindByLandCommodityID mocks base method.
func (m *MockHarvestRepository) FindByLandCommodityID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByLandCommodityID", ctx, id)
	ret0, _ := ret[0].([]*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByLandCommodityID indicates an expected call of FindByLandCommodityID.
func (mr *MockHarvestRepositoryMockRecorder) FindByLandCommodityID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByLandCommodityID", reflect.TypeOf((*MockHarvestRepository)(nil).FindByLandCommodityID), ctx, id)
}

// FindByLandID mocks base method.
func (m *MockHarvestRepository) FindByLandID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByLandID", ctx, id)
	ret0, _ := ret[0].([]*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByLandID indicates an expected call of FindByLandID.
func (mr *MockHarvestRepositoryMockRecorder) FindByLandID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByLandID", reflect.TypeOf((*MockHarvestRepository)(nil).FindByLandID), ctx, id)
}

// FindByRegionID mocks base method.
func (m *MockHarvestRepository) FindByRegionID(ctx context.Context, id uuid.UUID) ([]*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByRegionID", ctx, id)
	ret0, _ := ret[0].([]*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByRegionID indicates an expected call of FindByRegionID.
func (mr *MockHarvestRepositoryMockRecorder) FindByRegionID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByRegionID", reflect.TypeOf((*MockHarvestRepository)(nil).FindByRegionID), ctx, id)
}

// FindDeletedByID mocks base method.
func (m *MockHarvestRepository) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Harvest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeletedByID", ctx, id)
	ret0, _ := ret[0].(*domain.Harvest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeletedByID indicates an expected call of FindDeletedByID.
func (mr *MockHarvestRepositoryMockRecorder) FindDeletedByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeletedByID", reflect.TypeOf((*MockHarvestRepository)(nil).FindDeletedByID), ctx, id)
}

// Restore mocks base method.
func (m *MockHarvestRepository) Restore(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockHarvestRepositoryMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockHarvestRepository)(nil).Restore), ctx, id)
}

// Update mocks base method.
func (m *MockHarvestRepository) Update(ctx context.Context, id uuid.UUID, harvest *domain.Harvest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, harvest)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockHarvestRepositoryMockRecorder) Update(ctx, id, harvest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockHarvestRepository)(nil).Update), ctx, id, harvest)
}

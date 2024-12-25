// Code generated by MockGen. DO NOT EDIT.
// Source: service_api/repository/interface/sale_repository_interface.go

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

// MockSaleRepository is a mock of SaleRepository interface.
type MockSaleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSaleRepositoryMockRecorder
}

// MockSaleRepositoryMockRecorder is the mock recorder for MockSaleRepository.
type MockSaleRepositoryMockRecorder struct {
	mock *MockSaleRepository
}

// NewMockSaleRepository creates a new mock instance.
func NewMockSaleRepository(ctrl *gomock.Controller) *MockSaleRepository {
	mock := &MockSaleRepository{ctrl: ctrl}
	mock.recorder = &MockSaleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSaleRepository) EXPECT() *MockSaleRepositoryMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockSaleRepository) Count(ctx context.Context, filter *dto.PaginationFilterDTO) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, filter)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockSaleRepositoryMockRecorder) Count(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockSaleRepository)(nil).Count), ctx, filter)
}

// Create mocks base method.
func (m *MockSaleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, sale)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSaleRepositoryMockRecorder) Create(ctx, sale interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSaleRepository)(nil).Create), ctx, sale)
}

// Delete mocks base method.
func (m *MockSaleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSaleRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSaleRepository)(nil).Delete), ctx, id)
}

// DeletedCount mocks base method.
func (m *MockSaleRepository) DeletedCount(ctx context.Context, filter *dto.PaginationFilterDTO) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletedCount", ctx, filter)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletedCount indicates an expected call of DeletedCount.
func (mr *MockSaleRepositoryMockRecorder) DeletedCount(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletedCount", reflect.TypeOf((*MockSaleRepository)(nil).DeletedCount), ctx, filter)
}

// FindAll mocks base method.
func (m *MockSaleRepository) FindAll(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, params)
	ret0, _ := ret[0].([]*domain.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockSaleRepositoryMockRecorder) FindAll(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockSaleRepository)(nil).FindAll), ctx, params)
}

// FindAllDeleted mocks base method.
func (m *MockSaleRepository) FindAllDeleted(ctx context.Context, params *dto.PaginationDTO) ([]*domain.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllDeleted", ctx, params)
	ret0, _ := ret[0].([]*domain.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllDeleted indicates an expected call of FindAllDeleted.
func (mr *MockSaleRepositoryMockRecorder) FindAllDeleted(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllDeleted", reflect.TypeOf((*MockSaleRepository)(nil).FindAllDeleted), ctx, params)
}

// FindByCityID mocks base method.
func (m *MockSaleRepository) FindByCityID(ctx context.Context, params *dto.PaginationDTO, id int64) ([]*domain.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCityID", ctx, params, id)
	ret0, _ := ret[0].([]*domain.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCityID indicates an expected call of FindByCityID.
func (mr *MockSaleRepositoryMockRecorder) FindByCityID(ctx, params, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCityID", reflect.TypeOf((*MockSaleRepository)(nil).FindByCityID), ctx, params, id)
}

// FindByCommodityID mocks base method.
func (m *MockSaleRepository) FindByCommodityID(ctx context.Context, params *dto.PaginationDTO, id uuid.UUID) ([]*domain.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCommodityID", ctx, params, id)
	ret0, _ := ret[0].([]*domain.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCommodityID indicates an expected call of FindByCommodityID.
func (mr *MockSaleRepositoryMockRecorder) FindByCommodityID(ctx, params, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCommodityID", reflect.TypeOf((*MockSaleRepository)(nil).FindByCommodityID), ctx, params, id)
}

// FindByID mocks base method.
func (m *MockSaleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockSaleRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockSaleRepository)(nil).FindByID), ctx, id)
}

// FindDeletedByID mocks base method.
func (m *MockSaleRepository) FindDeletedByID(ctx context.Context, id uuid.UUID) (*domain.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDeletedByID", ctx, id)
	ret0, _ := ret[0].(*domain.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDeletedByID indicates an expected call of FindDeletedByID.
func (mr *MockSaleRepositoryMockRecorder) FindDeletedByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDeletedByID", reflect.TypeOf((*MockSaleRepository)(nil).FindDeletedByID), ctx, id)
}

// Restore mocks base method.
func (m *MockSaleRepository) Restore(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockSaleRepositoryMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockSaleRepository)(nil).Restore), ctx, id)
}

// Update mocks base method.
func (m *MockSaleRepository) Update(ctx context.Context, id uuid.UUID, sale *domain.Sale) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, sale)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSaleRepositoryMockRecorder) Update(ctx, id, sale interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSaleRepository)(nil).Update), ctx, id, sale)
}

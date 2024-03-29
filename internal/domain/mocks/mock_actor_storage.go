// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/actor_storage.go
//
// Generated by this command:
//
//	mockgen -source=internal/domain/actor_storage.go -destination=internal/domain/mocks/mock_actor_storage.go
//

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	reflect "reflect"

	domain "github.com/vaberof/vk-internship-task/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockActorStorage is a mock of ActorStorage interface.
type MockActorStorage struct {
	ctrl     *gomock.Controller
	recorder *MockActorStorageMockRecorder
}

// MockActorStorageMockRecorder is the mock recorder for MockActorStorage.
type MockActorStorageMockRecorder struct {
	mock *MockActorStorage
}

// NewMockActorStorage creates a new mock instance.
func NewMockActorStorage(ctrl *gomock.Controller) *MockActorStorage {
	mock := &MockActorStorage{ctrl: ctrl}
	mock.recorder = &MockActorStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActorStorage) EXPECT() *MockActorStorageMockRecorder {
	return m.recorder
}

// AreExists mocks base method.
func (m *MockActorStorage) AreExists(ids []domain.ActorId) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AreExists", ids)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AreExists indicates an expected call of AreExists.
func (mr *MockActorStorageMockRecorder) AreExists(ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AreExists", reflect.TypeOf((*MockActorStorage)(nil).AreExists), ids)
}

// Create mocks base method.
func (m *MockActorStorage) Create(name domain.ActorName, sex domain.ActorSex, birthDate domain.ActorBirthDate) (*domain.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name, sex, birthDate)
	ret0, _ := ret[0].(*domain.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockActorStorageMockRecorder) Create(name, sex, birthDate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockActorStorage)(nil).Create), name, sex, birthDate)
}

// Delete mocks base method.
func (m *MockActorStorage) Delete(id domain.ActorId) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockActorStorageMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockActorStorage)(nil).Delete), id)
}

// IsExists mocks base method.
func (m *MockActorStorage) IsExists(id domain.ActorId) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExists", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExists indicates an expected call of IsExists.
func (mr *MockActorStorageMockRecorder) IsExists(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExists", reflect.TypeOf((*MockActorStorage)(nil).IsExists), id)
}

// List mocks base method.
func (m *MockActorStorage) List(limit, offset int) ([]*domain.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", limit, offset)
	ret0, _ := ret[0].([]*domain.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockActorStorageMockRecorder) List(limit, offset any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockActorStorage)(nil).List), limit, offset)
}

// Update mocks base method.
func (m *MockActorStorage) Update(id domain.ActorId, name *domain.ActorName, sex *domain.ActorSex, birthDate *domain.ActorBirthDate) (*domain.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, name, sex, birthDate)
	ret0, _ := ret[0].(*domain.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockActorStorageMockRecorder) Update(id, name, sex, birthDate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockActorStorage)(nil).Update), id, name, sex, birthDate)
}

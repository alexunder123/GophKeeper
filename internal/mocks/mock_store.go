// Code generated by MockGen. DO NOT EDIT.
// Source: gophkeeper/internal/server/storage (interfaces: Storager)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// AuthUser mocks base method.
func (m *MockStorager) AuthUser(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthUser", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthUser indicates an expected call of AuthUser.
func (mr *MockStoragerMockRecorder) AuthUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUser", reflect.TypeOf((*MockStorager)(nil).AuthUser), arg0, arg1)
}

// ChangeUserPassword mocks base method.
func (m *MockStorager) ChangeUserPassword(arg0, arg1, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeUserPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeUserPassword indicates an expected call of ChangeUserPassword.
func (mr *MockStoragerMockRecorder) ChangeUserPassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeUserPassword", reflect.TypeOf((*MockStorager)(nil).ChangeUserPassword), arg0, arg1, arg2)
}

// CheckUser mocks base method.
func (m *MockStorager) CheckUser(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockStoragerMockRecorder) CheckUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockStorager)(nil).CheckUser), arg0)
}

// CloseDB mocks base method.
func (m *MockStorager) CloseDB() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseDB")
}

// CloseDB indicates an expected call of CloseDB.
func (mr *MockStoragerMockRecorder) CloseDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseDB", reflect.TypeOf((*MockStorager)(nil).CloseDB))
}

// RegisterUser mocks base method.
func (m *MockStorager) RegisterUser(arg0, arg1 string) (string, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockStoragerMockRecorder) RegisterUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockStorager)(nil).RegisterUser), arg0, arg1)
}

// UpdateUserData mocks base method.
func (m *MockStorager) UpdateUserData(arg0, arg1, arg2 string, arg3 []byte) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserData", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateUserData indicates an expected call of UpdateUserData.
func (mr *MockStoragerMockRecorder) UpdateUserData(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserData", reflect.TypeOf((*MockStorager)(nil).UpdateUserData), arg0, arg1, arg2, arg3)
}

// UsersData mocks base method.
func (m *MockStorager) UsersData(arg0 string) ([]byte, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersData", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// UsersData indicates an expected call of UsersData.
func (mr *MockStoragerMockRecorder) UsersData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersData", reflect.TypeOf((*MockStorager)(nil).UsersData), arg0)
}

// UsersDataLock mocks base method.
func (m *MockStorager) UsersDataLock(arg0, arg1 string) (bool, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersDataLock", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// UsersDataLock indicates an expected call of UsersDataLock.
func (mr *MockStoragerMockRecorder) UsersDataLock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersDataLock", reflect.TypeOf((*MockStorager)(nil).UsersDataLock), arg0, arg1)
}

// UsersTimeStamp mocks base method.
func (m *MockStorager) UsersTimeStamp(arg0 string) (string, bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersTimeStamp", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// UsersTimeStamp indicates an expected call of UsersTimeStamp.
func (mr *MockStoragerMockRecorder) UsersTimeStamp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersTimeStamp", reflect.TypeOf((*MockStorager)(nil).UsersTimeStamp), arg0)
}
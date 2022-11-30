// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/postgres/postgres.go

// Package mock_postgres is a generated GoMock package.
package postgres

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPostgres is a mock of Postgres interface.
type MockPostgres struct {
	ctrl     *gomock.Controller
	recorder *MockPostgresMockRecorder
}

// MockPostgresMockRecorder is the mock recorder for MockPostgres.
type MockPostgresMockRecorder struct {
	mock *MockPostgres
}

// NewMockPostgres creates a new mock instance.
func NewMockPostgres(ctrl *gomock.Controller) *MockPostgres {
	mock := &MockPostgres{ctrl: ctrl}
	mock.recorder = &MockPostgresMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostgres) EXPECT() *MockPostgresMockRecorder {
	return m.recorder
}

// GetAllURLs mocks base method.
func (m *MockPostgres) GetAllURLs() (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllURLs")
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllURLs indicates an expected call of GetAllURLs.
func (mr *MockPostgresMockRecorder) GetAllURLs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllURLs", reflect.TypeOf((*MockPostgres)(nil).GetAllURLs))
}

// GetOriginalURLByShortURL mocks base method.
func (m *MockPostgres) GetOriginalURLByShortURL(shortURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURLByShortURL", shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURLByShortURL indicates an expected call of GetOriginalURLByShortURL.
func (mr *MockPostgresMockRecorder) GetOriginalURLByShortURL(shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURLByShortURL", reflect.TypeOf((*MockPostgres)(nil).GetOriginalURLByShortURL), shortURL)
}

// GetShortURLByOriginalURL mocks base method.
func (m *MockPostgres) GetShortURLByOriginalURL(originalURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURLByOriginalURL", originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURLByOriginalURL indicates an expected call of GetShortURLByOriginalURL.
func (mr *MockPostgresMockRecorder) GetShortURLByOriginalURL(originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURLByOriginalURL", reflect.TypeOf((*MockPostgres)(nil).GetShortURLByOriginalURL), originalURL)
}

// Ping mocks base method.
func (m *MockPostgres) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPostgresMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPostgres)(nil).Ping), ctx)
}

// StoreURL mocks base method.
func (m *MockPostgres) StoreURL(originalURL, shortURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreURL", originalURL, shortURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreURL indicates an expected call of StoreURL.
func (mr *MockPostgresMockRecorder) StoreURL(originalURL, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreURL", reflect.TypeOf((*MockPostgres)(nil).StoreURL), originalURL, shortURL)
}

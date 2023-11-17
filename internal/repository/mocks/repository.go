// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go
//
// Generated by this command:
//
//	mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
//
// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	repository "github.com/NRKA/HTTP-Server/internal/repository"
	reflect "reflect"

	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	gomock "go.uber.org/mock/gomock"
)

// MockArticleInterface is a mock of ArticleInterface interface.
type MockArticleInterface struct {
	ctrl     *gomock.Controller
	recorder *MockArticleInterfaceMockRecorder
}

// MockArticleInterfaceMockRecorder is the mock recorder for MockArticleInterface.
type MockArticleInterfaceMockRecorder struct {
	mock *MockArticleInterface
}

// NewMockArticleInterface creates a new mock instance.
func NewMockArticleInterface(ctrl *gomock.Controller) *MockArticleInterface {
	mock := &MockArticleInterface{ctrl: ctrl}
	mock.recorder = &MockArticleInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleInterface) EXPECT() *MockArticleInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockArticleInterface) Create(ctx context.Context, article repository.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, article)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockArticleInterfaceMockRecorder) Create(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArticleInterface)(nil).Create), ctx, article)
}

// Delete mocks base method.
func (m *MockArticleInterface) Delete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockArticleInterfaceMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockArticleInterface)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockArticleInterface) GetByID(ctx context.Context, id int64) (repository.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(repository.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockArticleInterfaceMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockArticleInterface)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *MockArticleInterface) Update(ctx context.Context, article repository.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, article)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockArticleInterfaceMockRecorder) Update(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockArticleInterface)(nil).Update), ctx, article)
}

// MockDataBaseInterface is a mock of DataBaseInterface interface.
type MockDataBaseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDataBaseInterfaceMockRecorder
}

// MockDataBaseInterfaceMockRecorder is the mock recorder for MockDataBaseInterface.
type MockDataBaseInterfaceMockRecorder struct {
	mock *MockDataBaseInterface
}

// NewMockDataBaseInterface creates a new mock instance.
func NewMockDataBaseInterface(ctrl *gomock.Controller) *MockDataBaseInterface {
	mock := &MockDataBaseInterface{ctrl: ctrl}
	mock.recorder = &MockDataBaseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataBaseInterface) EXPECT() *MockDataBaseInterfaceMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockDataBaseInterface) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockDataBaseInterfaceMockRecorder) Exec(ctx, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockDataBaseInterface)(nil).Exec), varargs...)
}

// ExecQueryRow mocks base method.
func (m *MockDataBaseInterface) ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	m.ctrl.T.Helper()
	varargs := []any{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecQueryRow", varargs...)
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// ExecQueryRow indicates an expected call of ExecQueryRow.
func (mr *MockDataBaseInterfaceMockRecorder) ExecQueryRow(ctx, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecQueryRow", reflect.TypeOf((*MockDataBaseInterface)(nil).ExecQueryRow), varargs...)
}

// Get mocks base method.
func (m *MockDataBaseInterface) Get(ctx context.Context, dest any, query string, args ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockDataBaseInterfaceMockRecorder) Get(ctx, dest, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDataBaseInterface)(nil).Get), varargs...)
}

// GetPool mocks base method.
func (m *MockDataBaseInterface) GetPool() *pgxpool.Pool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPool")
	ret0, _ := ret[0].(*pgxpool.Pool)
	return ret0
}

// GetPool indicates an expected call of GetPool.
func (mr *MockDataBaseInterfaceMockRecorder) GetPool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPool", reflect.TypeOf((*MockDataBaseInterface)(nil).GetPool))
}

// Select mocks base method.
func (m *MockDataBaseInterface) Select(ctx context.Context, dest any, query string, args ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Select", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Select indicates an expected call of Select.
func (mr *MockDataBaseInterfaceMockRecorder) Select(ctx, dest, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockDataBaseInterface)(nil).Select), varargs...)
}

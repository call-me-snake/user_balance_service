// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/model.go

// Package mock_model is a generated GoMock package.
package mock_model

import (
	reflect "reflect"

	model "github.com/call-me-snake/user_balance_service/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockIBalanceInfoStorage is a mock of IBalanceInfoStorage interface.
type MockIBalanceInfoStorage struct {
	ctrl     *gomock.Controller
	recorder *MockIBalanceInfoStorageMockRecorder
}

// MockIBalanceInfoStorageMockRecorder is the mock recorder for MockIBalanceInfoStorage.
type MockIBalanceInfoStorageMockRecorder struct {
	mock *MockIBalanceInfoStorage
}

// NewMockIBalanceInfoStorage creates a new mock instance.
func NewMockIBalanceInfoStorage(ctrl *gomock.Controller) *MockIBalanceInfoStorage {
	mock := &MockIBalanceInfoStorage{ctrl: ctrl}
	mock.recorder = &MockIBalanceInfoStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBalanceInfoStorage) EXPECT() *MockIBalanceInfoStorageMockRecorder {
	return m.recorder
}

// GetAccountBalance mocks base method.
func (m *MockIBalanceInfoStorage) GetAccountBalance(id int) (*model.BalanceInfo, *model.CustomErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountBalance", id)
	ret0, _ := ret[0].(*model.BalanceInfo)
	ret1, _ := ret[1].(*model.CustomErr)
	return ret0, ret1
}

// GetAccountBalance indicates an expected call of GetAccountBalance.
func (mr *MockIBalanceInfoStorageMockRecorder) GetAccountBalance(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountBalance", reflect.TypeOf((*MockIBalanceInfoStorage)(nil).GetAccountBalance), id)
}

// ChangeAccountBalance mocks base method.
func (m *MockIBalanceInfoStorage) ChangeAccountBalance(id int, delta float64) (string, *model.CustomErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAccountBalance", id, delta)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*model.CustomErr)
	return ret0, ret1
}

// ChangeAccountBalance indicates an expected call of ChangeAccountBalance.
func (mr *MockIBalanceInfoStorageMockRecorder) ChangeAccountBalance(id, delta interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAccountBalance", reflect.TypeOf((*MockIBalanceInfoStorage)(nil).ChangeAccountBalance), id, delta)
}

// TransferSumBetweenAccounts mocks base method.
func (m *MockIBalanceInfoStorage) TransferSumBetweenAccounts(id1, id2 int, delta float64) (string, *model.CustomErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferSumBetweenAccounts", id1, id2, delta)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*model.CustomErr)
	return ret0, ret1
}

// TransferSumBetweenAccounts indicates an expected call of TransferSumBetweenAccounts.
func (mr *MockIBalanceInfoStorageMockRecorder) TransferSumBetweenAccounts(id1, id2, delta interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferSumBetweenAccounts", reflect.TypeOf((*MockIBalanceInfoStorage)(nil).TransferSumBetweenAccounts), id1, id2, delta)
}

// GetSortedTransactionsHistory mocks base method.
func (m *MockIBalanceInfoStorage) GetSortedTransactionsHistory(id int, sortedBy string, sortedByDesc bool) ([]model.TransactionRecord, *model.CustomErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSortedTransactionsHistory", id, sortedBy, sortedByDesc)
	ret0, _ := ret[0].([]model.TransactionRecord)
	ret1, _ := ret[1].(*model.CustomErr)
	return ret0, ret1
}

// GetSortedTransactionsHistory indicates an expected call of GetSortedTransactionsHistory.
func (mr *MockIBalanceInfoStorageMockRecorder) GetSortedTransactionsHistory(id, sortedBy, sortedByDesc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSortedTransactionsHistory", reflect.TypeOf((*MockIBalanceInfoStorage)(nil).GetSortedTransactionsHistory), id, sortedBy, sortedByDesc)
}

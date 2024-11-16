// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/repository/interfaces/payment.go

// Package mockRepository is a generated GoMock package.
package mockRepository

import (
	reflect "reflect"

	domain "github.com/ShahabazSulthan/Friendzy_Auth/pkg/domain"
	responsemodels "github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
	gomock "github.com/golang/mock/gomock"
)

// MockIPaymentRepository is a mock of IPaymentRepository interface.
type MockIPaymentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIPaymentRepositoryMockRecorder
}

// MockIPaymentRepositoryMockRecorder is the mock recorder for MockIPaymentRepository.
type MockIPaymentRepositoryMockRecorder struct {
	mock *MockIPaymentRepository
}

// NewMockIPaymentRepository creates a new mock instance.
func NewMockIPaymentRepository(ctrl *gomock.Controller) *MockIPaymentRepository {
	mock := &MockIPaymentRepository{ctrl: ctrl}
	mock.recorder = &MockIPaymentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPaymentRepository) EXPECT() *MockIPaymentRepositoryMockRecorder {
	return m.recorder
}

// CreateBlueTickVerification mocks base method.
func (m *MockIPaymentRepository) CreateBlueTickVerification(userID uint, verificationID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBlueTickVerification", userID, verificationID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBlueTickVerification indicates an expected call of CreateBlueTickVerification.
func (mr *MockIPaymentRepositoryMockRecorder) CreateBlueTickVerification(userID, verificationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBlueTickVerification", reflect.TypeOf((*MockIPaymentRepository)(nil).CreateBlueTickVerification), userID, verificationID)
}

// GetAllVerifiedUsers mocks base method.
func (m *MockIPaymentRepository) GetAllVerifiedUsers(limit, offset int) (*[]responsemodels.BlueTickResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllVerifiedUsers", limit, offset)
	ret0, _ := ret[0].(*[]responsemodels.BlueTickResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllVerifiedUsers indicates an expected call of GetAllVerifiedUsers.
func (mr *MockIPaymentRepositoryMockRecorder) GetAllVerifiedUsers(limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllVerifiedUsers", reflect.TypeOf((*MockIPaymentRepository)(nil).GetAllVerifiedUsers), limit, offset)
}

// GetBlueTickVerificationPrice mocks base method.
func (m *MockIPaymentRepository) GetBlueTickVerificationPrice() uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlueTickVerificationPrice")
	ret0, _ := ret[0].(uint)
	return ret0
}

// GetBlueTickVerificationPrice indicates an expected call of GetBlueTickVerificationPrice.
func (mr *MockIPaymentRepositoryMockRecorder) GetBlueTickVerificationPrice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlueTickVerificationPrice", reflect.TypeOf((*MockIPaymentRepository)(nil).GetBlueTickVerificationPrice))
}

// IsUserVerified mocks base method.
func (m *MockIPaymentRepository) IsUserVerified(userID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserVerified", userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserVerified indicates an expected call of IsUserVerified.
func (mr *MockIPaymentRepositoryMockRecorder) IsUserVerified(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserVerified", reflect.TypeOf((*MockIPaymentRepository)(nil).IsUserVerified), userID)
}

// OnlinePayment mocks base method.
func (m *MockIPaymentRepository) OnlinePayment(userID, verificationID string) (*responsemodels.OnlinePayment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnlinePayment", userID, verificationID)
	ret0, _ := ret[0].(*responsemodels.OnlinePayment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OnlinePayment indicates an expected call of OnlinePayment.
func (mr *MockIPaymentRepositoryMockRecorder) OnlinePayment(userID, verificationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnlinePayment", reflect.TypeOf((*MockIPaymentRepository)(nil).OnlinePayment), userID, verificationID)
}

// UpdateBlueTickPaymentSuccess mocks base method.
func (m *MockIPaymentRepository) UpdateBlueTickPaymentSuccess(verificationID string) (*domain.BlueTickVerification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBlueTickPaymentSuccess", verificationID)
	ret0, _ := ret[0].(*domain.BlueTickVerification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBlueTickPaymentSuccess indicates an expected call of UpdateBlueTickPaymentSuccess.
func (mr *MockIPaymentRepositoryMockRecorder) UpdateBlueTickPaymentSuccess(verificationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBlueTickPaymentSuccess", reflect.TypeOf((*MockIPaymentRepository)(nil).UpdateBlueTickPaymentSuccess), verificationID)
}

// UpdateBluetickStatus mocks base method.
func (m *MockIPaymentRepository) UpdateBluetickStatus(userID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBluetickStatus", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBluetickStatus indicates an expected call of UpdateBluetickStatus.
func (mr *MockIPaymentRepositoryMockRecorder) UpdateBluetickStatus(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBluetickStatus", reflect.TypeOf((*MockIPaymentRepository)(nil).UpdateBluetickStatus), userID)
}

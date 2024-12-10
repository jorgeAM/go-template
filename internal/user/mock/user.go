// Code generated by MockGen. DO NOT EDIT.
// Source: ./user.go
//
// Generated by this command:
//
//	mockgen -source=./user.go -destination=../mock/user.go -package=mock -mock_names=Repository=MockUserRepository
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	domain "github.com/jorgeAM/go-template/internal/user/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
	isgomock struct{}
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// FindByID mocks base method.
func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserRepositoryMockRecorder) FindByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserRepository)(nil).FindByID), ctx, id)
}

// Save mocks base method.
func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockUserRepositoryMockRecorder) Save(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserRepository)(nil).Save), ctx, user)
}
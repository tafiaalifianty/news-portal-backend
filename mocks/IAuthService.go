// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "final-project-backend/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// IAuthService is an autogenerated mock type for the IAuthService type
type IAuthService struct {
	mock.Mock
}

// Login provides a mock function with given fields: email, password
func (_m *IAuthService) Login(email string, password string) (*models.User, *string, *string, error) {
	ret := _m.Called(email, password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string, string) *models.User); ok {
		r0 = rf(email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 *string
	if rf, ok := ret.Get(1).(func(string, string) *string); ok {
		r1 = rf(email, password)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*string)
		}
	}

	var r2 *string
	if rf, ok := ret.Get(2).(func(string, string) *string); ok {
		r2 = rf(email, password)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*string)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string, string) error); ok {
		r3 = rf(email, password)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// RefreshAccessToken provides a mock function with given fields: userID, token
func (_m *IAuthService) RefreshAccessToken(userID int64, token *string) (*string, error) {
	ret := _m.Called(userID, token)

	var r0 *string
	if rf, ok := ret.Get(0).(func(int64, *string) *string); ok {
		r0 = rf(userID, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, *string) error); ok {
		r1 = rf(userID, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: user, referredCode
func (_m *IAuthService) Register(user *models.User, referredCode *string) (*models.User, *string, *string, error) {
	ret := _m.Called(user, referredCode)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(*models.User, *string) *models.User); ok {
		r0 = rf(user, referredCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 *string
	if rf, ok := ret.Get(1).(func(*models.User, *string) *string); ok {
		r1 = rf(user, referredCode)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*string)
		}
	}

	var r2 *string
	if rf, ok := ret.Get(2).(func(*models.User, *string) *string); ok {
		r2 = rf(user, referredCode)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*string)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(*models.User, *string) error); ok {
		r3 = rf(user, referredCode)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

type mockConstructorTestingTNewIAuthService interface {
	mock.TestingT
	Cleanup(func())
}

// NewIAuthService creates a new instance of IAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIAuthService(t mockConstructorTestingTNewIAuthService) *IAuthService {
	mock := &IAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

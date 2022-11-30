// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "final-project-backend/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// IUserSubscriptionService is an autogenerated mock type for the IUserSubscriptionService type
type IUserSubscriptionService struct {
	mock.Mock
}

// AddUserSubscription provides a mock function with given fields: userID, subscriptionID
func (_m *IUserSubscriptionService) AddUserSubscription(userID int64, subscriptionID int64) (*models.UserSubscriptions, error) {
	ret := _m.Called(userID, subscriptionID)

	var r0 *models.UserSubscriptions
	if rf, ok := ret.Get(0).(func(int64, int64) *models.UserSubscriptions); ok {
		r0 = rf(userID, subscriptionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserSubscriptions)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(userID, subscriptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllUserSubscriptions provides a mock function with given fields: userID
func (_m *IUserSubscriptionService) GetAllUserSubscriptions(userID int64) ([]*models.UserSubscriptions, error) {
	ret := _m.Called(userID)

	var r0 []*models.UserSubscriptions
	if rf, ok := ret.Get(0).(func(int64) []*models.UserSubscriptions); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.UserSubscriptions)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateUserQuota provides a mock function with given fields: userID, quotaNeeded
func (_m *IUserSubscriptionService) ValidateUserQuota(userID int64, quotaNeeded int) error {
	ret := _m.Called(userID, quotaNeeded)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int) error); ok {
		r0 = rf(userID, quotaNeeded)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewIUserSubscriptionService interface {
	mock.TestingT
	Cleanup(func())
}

// NewIUserSubscriptionService creates a new instance of IUserSubscriptionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIUserSubscriptionService(t mockConstructorTestingTNewIUserSubscriptionService) *IUserSubscriptionService {
	mock := &IUserSubscriptionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

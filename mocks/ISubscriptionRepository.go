// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "final-project-backend/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// ISubscriptionRepository is an autogenerated mock type for the ISubscriptionRepository type
type ISubscriptionRepository struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *ISubscriptionRepository) GetAll() ([]*models.Subscription, error) {
	ret := _m.Called()

	var r0 []*models.Subscription
	if rf, ok := ret.Get(0).(func() []*models.Subscription); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *ISubscriptionRepository) GetByID(id int64) (*models.Subscription, error) {
	ret := _m.Called(id)

	var r0 *models.Subscription
	if rf, ok := ret.Get(0).(func(int64) *models.Subscription); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewISubscriptionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewISubscriptionRepository creates a new instance of ISubscriptionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewISubscriptionRepository(t mockConstructorTestingTNewISubscriptionRepository) *ISubscriptionRepository {
	mock := &ISubscriptionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
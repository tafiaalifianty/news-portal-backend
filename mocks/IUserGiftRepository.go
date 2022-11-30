// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "final-project-backend/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// IUserGiftRepository is an autogenerated mock type for the IUserGiftRepository type
type IUserGiftRepository struct {
	mock.Mock
}

// Insert provides a mock function with given fields: userGift
func (_m *IUserGiftRepository) Insert(userGift *models.UserGift) (*models.UserGift, error) {
	ret := _m.Called(userGift)

	var r0 *models.UserGift
	if rf, ok := ret.Get(0).(func(*models.UserGift) *models.UserGift); ok {
		r0 = rf(userGift)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserGift)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.UserGift) error); ok {
		r1 = rf(userGift)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIUserGiftRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIUserGiftRepository creates a new instance of IUserGiftRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIUserGiftRepository(t mockConstructorTestingTNewIUserGiftRepository) *IUserGiftRepository {
	mock := &IUserGiftRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

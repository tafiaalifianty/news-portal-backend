// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "final-project-backend/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// IUserVoucherRepository is an autogenerated mock type for the IUserVoucherRepository type
type IUserVoucherRepository struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *IUserVoucherRepository) GetAll() ([]*models.UserVoucher, error) {
	ret := _m.Called()

	var r0 []*models.UserVoucher
	if rf, ok := ret.Get(0).(func() []*models.UserVoucher); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.UserVoucher)
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

// GetAllByUserID provides a mock function with given fields: userID
func (_m *IUserVoucherRepository) GetAllByUserID(userID int64) ([]*models.UserVoucher, error) {
	ret := _m.Called(userID)

	var r0 []*models.UserVoucher
	if rf, ok := ret.Get(0).(func(int64) []*models.UserVoucher); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.UserVoucher)
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

// GetByCode provides a mock function with given fields: code
func (_m *IUserVoucherRepository) GetByCode(code string) (*models.UserVoucher, error) {
	ret := _m.Called(code)

	var r0 *models.UserVoucher
	if rf, ok := ret.Get(0).(func(string) *models.UserVoucher); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserVoucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: userVoucher
func (_m *IUserVoucherRepository) Insert(userVoucher *models.UserVoucher) (*models.UserVoucher, error) {
	ret := _m.Called(userVoucher)

	var r0 *models.UserVoucher
	if rf, ok := ret.Get(0).(func(*models.UserVoucher) *models.UserVoucher); ok {
		r0 = rf(userVoucher)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserVoucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.UserVoucher) error); ok {
		r1 = rf(userVoucher)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: userVoucher
func (_m *IUserVoucherRepository) Update(userVoucher *models.UserVoucher) (*models.UserVoucher, int, error) {
	ret := _m.Called(userVoucher)

	var r0 *models.UserVoucher
	if rf, ok := ret.Get(0).(func(*models.UserVoucher) *models.UserVoucher); ok {
		r0 = rf(userVoucher)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserVoucher)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(*models.UserVoucher) int); ok {
		r1 = rf(userVoucher)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*models.UserVoucher) error); ok {
		r2 = rf(userVoucher)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewIUserVoucherRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIUserVoucherRepository creates a new instance of IUserVoucherRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIUserVoucherRepository(t mockConstructorTestingTNewIUserVoucherRepository) *IUserVoucherRepository {
	mock := &IUserVoucherRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

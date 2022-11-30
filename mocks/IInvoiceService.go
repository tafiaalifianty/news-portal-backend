// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "final-project-backend/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// IInvoiceService is an autogenerated mock type for the IInvoiceService type
type IInvoiceService struct {
	mock.Mock
}

// Create provides a mock function with given fields: invoice
func (_m *IInvoiceService) Create(invoice *models.Invoice) (*models.Invoice, error) {
	ret := _m.Called(invoice)

	var r0 *models.Invoice
	if rf, ok := ret.Get(0).(func(*models.Invoice) *models.Invoice); ok {
		r0 = rf(invoice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Invoice) error); ok {
		r1 = rf(invoice)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *IInvoiceService) GetAll() ([]*models.Invoice, error) {
	ret := _m.Called()

	var r0 []*models.Invoice
	if rf, ok := ret.Get(0).(func() []*models.Invoice); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Invoice)
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

// GetByCode provides a mock function with given fields: code
func (_m *IInvoiceService) GetByCode(code string) (*models.Invoice, error) {
	ret := _m.Called(code)

	var r0 *models.Invoice
	if rf, ok := ret.Get(0).(func(string) *models.Invoice); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
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

// GetByID provides a mock function with given fields: id
func (_m *IInvoiceService) GetByID(id int64) (*models.Invoice, error) {
	ret := _m.Called(id)

	var r0 *models.Invoice
	if rf, ok := ret.Get(0).(func(int64) *models.Invoice); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
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

// GetUserInvoiceByCode provides a mock function with given fields: code, userID
func (_m *IInvoiceService) GetUserInvoiceByCode(code string, userID int64) (*models.Invoice, error) {
	ret := _m.Called(code, userID)

	var r0 *models.Invoice
	if rf, ok := ret.Get(0).(func(string, int64) *models.Invoice); ok {
		r0 = rf(code, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int64) error); ok {
		r1 = rf(code, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserInvoices provides a mock function with given fields: userID
func (_m *IInvoiceService) GetUserInvoices(userID int64) ([]*models.Invoice, error) {
	ret := _m.Called(userID)

	var r0 []*models.Invoice
	if rf, ok := ret.Get(0).(func(int64) []*models.Invoice); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Invoice)
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

// UpdateStatus provides a mock function with given fields: code, status
func (_m *IInvoiceService) UpdateStatus(code string, status models.InvoiceStatus) (*models.Invoice, []*models.Gift, []*models.Voucher, error) {
	ret := _m.Called(code, status)

	var r0 *models.Invoice
	if rf, ok := ret.Get(0).(func(string, models.InvoiceStatus) *models.Invoice); ok {
		r0 = rf(code, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invoice)
		}
	}

	var r1 []*models.Gift
	if rf, ok := ret.Get(1).(func(string, models.InvoiceStatus) []*models.Gift); ok {
		r1 = rf(code, status)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*models.Gift)
		}
	}

	var r2 []*models.Voucher
	if rf, ok := ret.Get(2).(func(string, models.InvoiceStatus) []*models.Voucher); ok {
		r2 = rf(code, status)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).([]*models.Voucher)
		}
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(string, models.InvoiceStatus) error); ok {
		r3 = rf(code, status)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

type mockConstructorTestingTNewIInvoiceService interface {
	mock.TestingT
	Cleanup(func())
}

// NewIInvoiceService creates a new instance of IInvoiceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIInvoiceService(t mockConstructorTestingTNewIInvoiceService) *IInvoiceService {
	mock := &IInvoiceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
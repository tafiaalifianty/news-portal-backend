package services

import (
	"fmt"
	"testing"

	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/models"
	"final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNewInvoiceService(t *testing.T) {
	NewInvoiceService(&InvoiceServiceConfig{
		invoiceRepository: mocks.NewIInvoiceRepository(t),
	})
}

func Test_invoiceService_GetAll(t *testing.T) {
	mockInvoices := []*models.Invoice{}
	mockError := fmt.Errorf("error")
	type fields struct {
		invoiceRepository *mocks.IInvoiceRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.IInvoiceRepository)
		want        []*models.Invoice
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from invoiceRepository.GetAll",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetAll", &models.Invoice{}).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetAll", &models.Invoice{}).Return(mockInvoices, nil)
			},
			want:        mockInvoices,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &invoiceService{
				invoiceRepository: tt.fields.invoiceRepository,
			}

			tt.mock(tt.fields.invoiceRepository)
			got, err := s.GetAll()

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_invoiceService_Create(t *testing.T) {
	mockInvoice := &models.Invoice{}

	mockError := fmt.Errorf("error")
	type fields struct {
		invoiceRepository *mocks.IInvoiceRepository
	}
	type args struct {
		invoice *models.Invoice
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IInvoiceRepository)
		want        *models.Invoice
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from invoiceRepository.Insert",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			args: args{
				invoice: mockInvoice,
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("Insert", mockInvoice).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			args: args{
				invoice: mockInvoice,
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("Insert", mockInvoice).Return(mockInvoice, nil)
			},
			want:        mockInvoice,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &invoiceService{
				invoiceRepository: tt.fields.invoiceRepository,
			}

			tt.mock(tt.fields.invoiceRepository)
			got, err := s.Create(tt.args.invoice)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_invoiceService_GetByID(t *testing.T) {
	mockInvoice := &models.Invoice{}
	mockID := int64(1)
	mockError := fmt.Errorf("error")
	type fields struct {
		invoiceRepository *mocks.IInvoiceRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.IInvoiceRepository)
		want        *models.Invoice
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from invoiceRepository.GetByID",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetByID", mockID).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetByID", mockID).Return(mockInvoice, nil)
			},
			want:        mockInvoice,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &invoiceService{
				invoiceRepository: tt.fields.invoiceRepository,
			}

			tt.mock(tt.fields.invoiceRepository)
			got, err := s.GetByID(mockID)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_invoiceService_UpdateStatus(t *testing.T) {
	mockID := int64(1)
	mockCode := "code"
	mockError := fmt.Errorf("error")
	// defaultMockInvoice := &models.Invoice{
	// 	Model: gorm.Model{
	// 		ID: uint(mockID),
	// 	},
	// 	Code:   mockCode,
	// 	Status: models.PROCESSED,
	// }
	// completedMockInvoice := &models.Invoice{
	// 	Model: gorm.Model{
	// 		ID: uint(mockID),
	// 	},
	// 	Code:           mockCode,
	// 	Status:         models.COMPLETED,
	// 	SubscriptionID: 1,
	// 	UserID:         1,
	// }
	// mockSubscription := &models.Subscription{
	// 	Quota: 1,
	// }
	defaultStatusParams := models.PROCESSED
	defaultMockReturnGetByID := &models.Invoice{
		Status: models.WAITING,
		Code:   mockCode,
		Model: gorm.Model{
			ID: uint(mockID),
		},
	}

	type args struct {
		code   string
		status models.InvoiceStatus
	}
	type fields struct {
		invoiceRepository          *mocks.IInvoiceRepository
		subscriptionRepository     *mocks.ISubscriptionRepository
		userSubscriptionRepository *mocks.IUserSubscriptionRepository
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IInvoiceRepository, *mocks.ISubscriptionRepository, *mocks.IUserSubscriptionRepository)
		want        *models.Invoice
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from invoiceRepository.GetByCode",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			args: args{
				code: mockCode,
			},
			mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
				r.On("GetByCode", mockCode).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "ERROR | Invalid Status Update (besides Rejected & not from Processed)",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			args: args{
				code:   mockCode,
				status: models.WAITING,
			},
			mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
				r.On("GetByCode", mockCode).Return(&models.Invoice{Status: models.WAITING}, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: errn.ErrInvalidStatusUpdate,
		},
		{
			name: "ERROR | Invalid Status Update (Rejected & not from Processed)",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			args: args{
				code:   mockCode,
				status: models.REJECTED,
			},
			mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
				r.On("GetByCode", mockCode).Return(&models.Invoice{Status: models.WAITING}, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: errn.ErrInvalidStatusUpdate,
		},
		{
			name: "ERROR | No data updated from invoiceRepository.Update",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			args: args{
				code:   mockCode,
				status: defaultStatusParams,
			},
			mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
				r.On("GetByCode", mockCode).Return(defaultMockReturnGetByID, nil)
				r.On("Update", mock.MatchedBy(func(i *models.Invoice) bool {
					return true
				})).Return(nil, 0, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "ERROR | Error from invoiceRepository.Update",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			args: args{
				code:   mockCode,
				status: defaultStatusParams,
			},
			mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
				r.On("GetByCode", mockCode).Return(defaultMockReturnGetByID, nil)
				r.On("Update", mock.MatchedBy(func(i interface{}) bool {
					invoice := i.(*models.Invoice)
					return invoice.ID == uint(mockID) &&
						invoice.Status == defaultStatusParams
				})).Return(nil, 1, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		// {
		// 	name: "ERROR | Update to Completed | Error from subscriptionRepository.GetByID",
		// 	fields: fields{
		// 		invoiceRepository:          mocks.NewIInvoiceRepository(t),
		// 		subscriptionRepository:     mocks.NewISubscriptionRepository(t),
		// 		userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
		// 	},
		// 	args: args{
		// 		code:   mockCode,
		// 		status: defaultStatusParams,
		// 	},
		// 	mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
		// 		r.On("GetByCode", mockCode).Return(defaultMockReturnGetByID, nil)
		// 		r.On("Update", mock.MatchedBy(func(i interface{}) bool {
		// 			invoice := i.(*models.Invoice)
		// 			return invoice.ID == uint(mockID) &&
		// 				invoice.Status == defaultStatusParams
		// 		})).Return(completedMockInvoice, 1, nil)
		// 		s.On("GetByID", completedMockInvoice.SubscriptionID).Return(nil, mockError)
		// 	},
		// 	want:        nil,
		// 	wantErr:     true,
		// 	expectedErr: mockError,
		// },
		// {
		// 	name: "ERROR | Update to Completed | Error from userSubscriptionRepository.Insert",
		// 	fields: fields{
		// 		invoiceRepository:          mocks.NewIInvoiceRepository(t),
		// 		subscriptionRepository:     mocks.NewISubscriptionRepository(t),
		// 		userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
		// 	},
		// 	args: args{
		// 		code:   mockCode,
		// 		status: defaultStatusParams,
		// 	},
		// 	mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
		// 		r.On("GetByCode", mockCode).Return(defaultMockReturnGetByID, nil)
		// 		r.On("Update", mock.MatchedBy(func(i interface{}) bool {
		// 			invoice := i.(*models.Invoice)
		// 			return invoice.ID == uint(mockID) &&
		// 				invoice.Status == defaultStatusParams
		// 		})).Return(completedMockInvoice, 1, nil)
		// 		s.On("GetByID", completedMockInvoice.SubscriptionID).Return(mockSubscription, nil)
		// 		us.On("Insert", mock.MatchedBy(func(i interface{}) bool {
		// 			userSubscription := i.(*models.UserSubscriptions)
		// 			return userSubscription.UserID == completedMockInvoice.UserID &&
		// 				userSubscription.SubscriptionID == completedMockInvoice.SubscriptionID
		// 		})).Return(nil, mockError)
		// 	},
		// 	want:        nil,
		// 	wantErr:     true,
		// 	expectedErr: mockError,
		// },
		// {
		// 	name: "SUCCESS",
		// 	fields: fields{
		// 		invoiceRepository: mocks.NewIInvoiceRepository(t),
		// 	},
		// 	args: args{
		// 		code:   mockCode,
		// 		status: defaultStatusParams,
		// 	},
		// 	mock: func(r *mocks.IInvoiceRepository, s *mocks.ISubscriptionRepository, us *mocks.IUserSubscriptionRepository) {
		// 		r.On("GetByCode", mockCode).Return(defaultMockReturnGetByID, nil)
		// 		r.On("Update", mock.MatchedBy(func(i interface{}) bool {
		// 			invoice := i.(*models.Invoice)
		// 			return invoice.ID == uint(mockID) &&
		// 				invoice.Status == defaultStatusParams
		// 		})).Return(defaultMockInvoice, 1, nil)
		// 	},
		// 	want:        defaultMockInvoice,
		// 	wantErr:     false,
		// 	expectedErr: nil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &invoiceService{
				invoiceRepository:          tt.fields.invoiceRepository,
				subscriptionRepository:     tt.fields.subscriptionRepository,
				userSubscriptionRepository: tt.fields.userSubscriptionRepository,
			}

			tt.mock(tt.fields.invoiceRepository, tt.fields.subscriptionRepository, tt.fields.userSubscriptionRepository)
			got, _, _, err := s.UpdateStatus(tt.args.code, tt.args.status)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_invoiceService_GetUserInvoiceByCode(t *testing.T) {
	mockCode := "code"
	mockUserID := int64(1)
	mockInvoice := &models.Invoice{UserID: mockUserID}
	mockError := fmt.Errorf("error")
	type fields struct {
		invoiceRepository *mocks.IInvoiceRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.IInvoiceRepository)
		want        *models.Invoice
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from invoiceRepository.GetByCode",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetByCode", mockCode).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "ERROR | Error from different userID returned then the one provided",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetByCode", mockCode).Return(&models.Invoice{UserID: int64(2)}, nil)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: errn.ErrNotAuthorized,
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetByCode", mockCode).Return(mockInvoice, nil)
			},
			want:        mockInvoice,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &invoiceService{
				invoiceRepository: tt.fields.invoiceRepository,
			}

			tt.mock(tt.fields.invoiceRepository)
			got, err := s.GetUserInvoiceByCode(mockCode, mockUserID)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_invoiceService_GetByCode(t *testing.T) {
	mockInvoice := &models.Invoice{}
	mockCode := "code"
	mockError := fmt.Errorf("error")
	type fields struct {
		invoiceRepository *mocks.IInvoiceRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.IInvoiceRepository)
		want        *models.Invoice
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from invoiceRepository.GetByCode",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetByCode", mockCode).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceRepository: mocks.NewIInvoiceRepository(t),
			},
			mock: func(r *mocks.IInvoiceRepository) {
				r.On("GetByCode", mockCode).Return(mockInvoice, nil)
			},
			want:        mockInvoice,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &invoiceService{
				invoiceRepository: tt.fields.invoiceRepository,
			}

			tt.mock(tt.fields.invoiceRepository)
			got, err := s.GetByCode(mockCode)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"
	"final-project-backend/internal/services"
	"final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHandler_GetAllInvoice(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockInvoice := []*models.Invoice{}

	type fields struct {
		invoiceService mocks.IInvoiceService
	}
	tests := []struct {
		name   string
		fields fields
		mock   func(*mocks.IInvoiceService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Invoice.GetAll",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetAll").Return(nil, mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetAll").Return(mockInvoice, nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    []interface{}{},
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Invoice: &tt.fields.invoiceService,
				},
			}

			tt.mock(&tt.fields.invoiceService)
			r := helpers.SetUpRouter()

			endpoint := "/invoices"
			r.GET(endpoint, h.GetAllInvoices)

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_CreateInvoice(t *testing.T) {
	invalidRequest := &dtos.CreateInvoiceRequest{}
	validRequest := &dtos.CreateInvoiceRequest{
		OriginalPrice:  10000,
		SubscriptionID: 1,
	}

	mockValidDataInInterface, err := helpers.StructToMap(&dtos.CreateInvoiceResponse{})
	require.NoError(t, err)

	mockError := fmt.Errorf("error")
	type fields struct {
		invoiceService *mocks.IInvoiceService
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		mock                   func(*mocks.IInvoiceService)
		mockUserFromMiddleware bool
		want                   helpers.JsonResponse
	}{
		{
			name: "ERROR | Invalid Request Body",
			fields: fields{
				invoiceService: mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(invalidRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
			},
			mockUserFromMiddleware: false,
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Invalid User Context",
			fields: fields{
				invoiceService: mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
			},
			mockUserFromMiddleware: false,
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from InvoiceService.Create",
			fields: fields{
				invoiceService: mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("Create", &models.Invoice{
					UserID:         int64(1),
					OriginalPrice:  validRequest.OriginalPrice,
					Total:          validRequest.OriginalPrice,
					Status:         models.WAITING,
					SubscriptionID: int64(validRequest.SubscriptionID),
				}).Return(nil, mockError)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceService: mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("Create", &models.Invoice{
					UserID:         int64(1),
					OriginalPrice:  validRequest.OriginalPrice,
					Total:          validRequest.OriginalPrice,
					Status:         models.WAITING,
					SubscriptionID: int64(validRequest.SubscriptionID),
				}).Return(&models.Invoice{}, nil)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Invoice: tt.fields.invoiceService,
				},
			}

			tt.mock(tt.fields.invoiceService)
			r := helpers.SetUpRouter()
			endpoint := "/invoices"

			if tt.mockUserFromMiddleware {
				r.POST(endpoint, helpers.MiddlewareMockUser(dtos.JwtData{ID: int64(1)}), h.CreateInvoice)
			} else {
				r.POST(endpoint, h.CreateInvoice)
			}
			req, _ := http.NewRequest(
				http.MethodPost,
				endpoint,
				tt.args.body,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_GetWaitingInvoiceByCode(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockCode := "code"
	validResponse := &dtos.InvoiceResponse{
		Status: models.WAITING.String(),
	}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)

	type fields struct {
		invoiceService mocks.IInvoiceService
	}
	tests := []struct {
		name                     string
		fields                   fields
		mock                     func(*mocks.IInvoiceService)
		mockParamsFromMiddleware bool
		want                     helpers.JsonResponse
	}{
		{
			name: "ERROR | Error Invoice.GetByCode: gorm.ErrRecordNotFound",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetByCode", strings.ToUpper(mockCode)).Return(nil, gorm.ErrRecordNotFound)
			},
			mockParamsFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrNoInvoicesFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error Invoice.GetByCode: others",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetByCode", strings.ToUpper(mockCode)).Return(nil, mockError)
			},
			mockParamsFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Requested Invoice's status is not Waiting",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetByCode", strings.ToUpper(mockCode)).Return(&models.Invoice{
					Status: models.COMPLETED,
				}, nil)
			},
			mockParamsFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: errn.ErrInvoiceNotAwaitingPayment.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetByCode", strings.ToUpper(mockCode)).Return(&models.Invoice{
					Status: models.WAITING,
				}, nil)
			},
			mockParamsFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Invoice: &tt.fields.invoiceService,
				},
			}

			tt.mock(&tt.fields.invoiceService)
			r := helpers.SetUpRouter()

			endpoint := "/invoices"
			if tt.mockParamsFromMiddleware {
				r.GET(endpoint, helpers.MiddlewareMockParams("code", mockCode), h.GetWaitingInvoiceByCode)
			} else {
				r.GET(endpoint, h.GetWaitingInvoiceByCode)
			}

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func BoolPointer(b bool) *bool {
	return &b
}
func TestHandler_UpdateWaitingInvoiceToProcessed(t *testing.T) {
	mockCode := "code"
	mockError := fmt.Errorf("error")
	validResponse := &dtos.InvoiceResponse{}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)

	type fields struct {
		invoiceService mocks.IInvoiceService
	}
	tests := []struct {
		name   string
		fields fields
		mock   func(*mocks.IInvoiceService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Invoice.UpdateStatus: gorm.ErrRecordNotFound",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.PROCESSED).Return(nil, nil, nil, gorm.ErrRecordNotFound)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrNoInvoicesFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Invoice.UpdateStatus: ErrInvalidStatusUpdate",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.PROCESSED).Return(nil, nil, nil, errn.ErrInvalidStatusUpdate)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: errn.ErrInvalidStatusUpdate.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Invoice.UpdateStatus: Others",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.PROCESSED).Return(nil, nil, nil, mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.PROCESSED).Return(&models.Invoice{}, nil, nil, nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Invoice: &tt.fields.invoiceService,
				},
			}

			tt.mock(&tt.fields.invoiceService)
			r := helpers.SetUpRouter()

			endpoint := "/invoices"
			r.PATCH(endpoint, helpers.MiddlewareMockParams("code", mockCode), h.UpdateWaitingInvoiceToProcessed)

			req, _ := http.NewRequest(
				http.MethodPatch,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_UpdateProcessedInvoice(t *testing.T) {
	mockCode := "code"
	invalidRequest := dtos.UpdateProcessedInvoiceRequest{}
	validRequest := dtos.UpdateProcessedInvoiceRequest{
		IsSuccess: BoolPointer(true),
	}
	mockError := fmt.Errorf("error")
	validResponse := &dtos.UpdateProcessedInvoiceResponse{
		GiftsSent:    []*dtos.GiftResponse{},
		VouchersSent: []*dtos.VoucherResponse{},
	}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)

	type fields struct {
		invoiceService mocks.IInvoiceService
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		mock   func(*mocks.IInvoiceService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from invalid request body",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(invalidRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
			},

			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Invoice.UpdateStatus: gorm.ErrRecordNotFound",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.COMPLETED).Return(nil, nil, nil, gorm.ErrRecordNotFound)
			},

			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrNoInvoicesFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Invoice.UpdateStatus: errn.ErrInvalidStatusUpdate",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.COMPLETED).Return(nil, nil, nil, errn.ErrInvalidStatusUpdate)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: errn.ErrInvalidStatusUpdate.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from Invoice.UpdateStatus: other errors",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.COMPLETED).Return(nil, nil, nil, mockError)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(dtos.UpdateProcessedInvoiceRequest{
					IsSuccess: BoolPointer(false),
				}),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("UpdateStatus", strings.ToUpper(mockCode), models.REJECTED).Return(&models.Invoice{}, []*models.Gift{}, []*models.Voucher{}, nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Invoice: &tt.fields.invoiceService,
				},
			}

			tt.mock(&tt.fields.invoiceService)
			r := helpers.SetUpRouter()

			endpoint := "/invoices"
			r.PATCH(endpoint, helpers.MiddlewareMockParams("code", mockCode), h.UpdateProcessedInvoice)

			req, _ := http.NewRequest(
				http.MethodPatch,
				endpoint,
				tt.args.body,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestHandler_GetWaitingInvoiceByCodeProtected(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockCode := "code"
	mockUserID := int64(1)
	validResponse := &dtos.InvoiceResponse{
		Status: models.WAITING.String(),
	}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)

	type fields struct {
		invoiceService mocks.IInvoiceService
	}
	tests := []struct {
		name                   string
		fields                 fields
		mock                   func(*mocks.IInvoiceService)
		mockUserFromMiddleware bool
		want                   helpers.JsonResponse
	}{
		{
			name: "ERROR | Invalid user context",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
			},
			mockUserFromMiddleware: false,
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error Invoice.GetUserInvoiceByCode: gorm.ErrRecordNotFound",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetUserInvoiceByCode", strings.ToUpper(mockCode), mockUserID).Return(nil, gorm.ErrRecordNotFound)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrNoInvoicesFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error Invoice.GetUserInvoiceByCode: errn.ErrNotAuthorized",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetUserInvoiceByCode", strings.ToUpper(mockCode), mockUserID).Return(nil, errn.ErrNotAuthorized)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error Invoice.GetUserInvoiceByCode: others",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetUserInvoiceByCode", strings.ToUpper(mockCode), mockUserID).Return(nil, mockError)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Requested Invoice's status is not Waiting",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetUserInvoiceByCode", strings.ToUpper(mockCode), mockUserID).Return(&models.Invoice{
					Status: models.COMPLETED,
				}, nil)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: errn.ErrInvoiceNotAwaitingPayment.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				invoiceService: *mocks.NewIInvoiceService(t),
			},
			mock: func(s *mocks.IInvoiceService) {
				s.On("GetUserInvoiceByCode", strings.ToUpper(mockCode), mockUserID).Return(&models.Invoice{
					Status: models.WAITING,
				}, nil)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusOK,
				Message: http.StatusText(http.StatusOK),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					Invoice: &tt.fields.invoiceService,
				},
			}

			tt.mock(&tt.fields.invoiceService)
			r := helpers.SetUpRouter()

			endpoint := "/invoices"
			if tt.mockUserFromMiddleware {
				r.GET(endpoint, helpers.MiddlewareMockParams("code", "code"), helpers.MiddlewareMockUser(dtos.JwtData{ID: mockUserID}), h.GetWaitingInvoiceByCodeProtected)
			} else {
				r.GET(endpoint, helpers.MiddlewareMockParams("code", "code"), h.GetWaitingInvoiceByCodeProtected)
			}

			req, _ := http.NewRequest(
				http.MethodGet,
				endpoint,
				nil,
			)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var response helpers.JsonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.want.Code, w.Code)
			assert.Equal(t, tt.want, response)
		})
	}
}

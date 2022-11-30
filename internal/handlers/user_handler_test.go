package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestHandler_GetUserProfile(t *testing.T) {
	mockUserContext := dtos.JwtData{
		ID:    1,
		Email: "email",
		Role:  string(models.Member),
		Token: "token",
	}
	validResponse := &dtos.GetUserProfileResponseDTO{}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)
	mockError := fmt.Errorf("error")

	type fields struct {
		userService mocks.IUserService
	}
	tests := []struct {
		name                   string
		fields                 fields
		mock                   func(*mocks.IUserService)
		mockUserFromMiddleware bool
		want                   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Context.Get",
			fields: fields{
				userService: *mocks.NewIUserService(t),
			},
			mock: func(as *mocks.IUserService) {
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
			name: "ERROR | Error from UserService.GetByID (user not found)",
			fields: fields{
				userService: *mocks.NewIUserService(t),
			},
			mock: func(s *mocks.IUserService) {
				s.On("GetByID", mockUserContext.ID).Return(nil, gorm.ErrRecordNotFound)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusNotFound,
				Message: errn.ErrUserNotFound.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from UserService.GetByID (other errors)",
			fields: fields{
				userService: *mocks.NewIUserService(t),
			},
			mock: func(s *mocks.IUserService) {
				s.On("GetByID", mockUserContext.ID).Return(nil, mockError)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: mockError.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				userService: *mocks.NewIUserService(t),
			},
			mock: func(s *mocks.IUserService) {
				s.On("GetByID", mockUserContext.ID).Return(&models.User{}, nil)
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
					User: &tt.fields.userService,
				},
			}

			tt.mock(&tt.fields.userService)
			r := helpers.SetUpRouter()
			endpoint := "/user/profile"

			if tt.mockUserFromMiddleware {
				r.GET(endpoint, helpers.MiddlewareMockUser(mockUserContext), h.GetUserProfile)
			} else {
				r.GET(endpoint, h.GetUserProfile)
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

func TestHandler_GetUserHistories(t *testing.T) {
	mockUserContext := dtos.JwtData{
		ID:    1,
		Email: "email",
		Role:  string(models.Member),
		Token: "token",
	}
	mockError := fmt.Errorf("error")

	type fields struct {
		historyService mocks.IHistoryService
	}
	tests := []struct {
		name                   string
		fields                 fields
		mock                   func(*mocks.IHistoryService)
		mockUserFromMiddleware bool
		want                   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Context.Get",
			fields: fields{
				historyService: *mocks.NewIHistoryService(t),
			},
			mock: func(s *mocks.IHistoryService) {
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
			name: "ERROR | Error from History.GetByUserID",
			fields: fields{
				historyService: *mocks.NewIHistoryService(t),
			},
			mock: func(s *mocks.IHistoryService) {
				s.On("GetByUserID", mockUserContext.ID).Return(nil, mockError)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: mockError.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				historyService: *mocks.NewIHistoryService(t),
			},
			mock: func(s *mocks.IHistoryService) {
				s.On("GetByUserID", mockUserContext.ID).Return([]*models.History{}, nil)
			},
			mockUserFromMiddleware: true,
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
					History: &tt.fields.historyService,
				},
			}

			tt.mock(&tt.fields.historyService)
			r := helpers.SetUpRouter()
			endpoint := "/user/history"

			if tt.mockUserFromMiddleware {
				r.GET(endpoint, helpers.MiddlewareMockUser(mockUserContext), h.GetUserHistories)
			} else {
				r.GET(endpoint, h.GetUserHistories)
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

func TestHandler_GetUserSubscriptions(t *testing.T) {
	mockUserContext := dtos.JwtData{
		ID:    1,
		Email: "email",
		Role:  string(models.Member),
		Token: "token",
	}
	mockError := fmt.Errorf("error")

	type fields struct {
		userSubscriptionService mocks.IUserSubscriptionService
	}
	tests := []struct {
		name                   string
		fields                 fields
		mock                   func(*mocks.IUserSubscriptionService)
		mockUserFromMiddleware bool
		want                   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Context.Get",
			fields: fields{
				userSubscriptionService: *mocks.NewIUserSubscriptionService(t),
			},
			mock: func(s *mocks.IUserSubscriptionService) {
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
			name: "ERROR | Error from UserSubscription.GetAllUserSubscriptions",
			fields: fields{
				userSubscriptionService: *mocks.NewIUserSubscriptionService(t),
			},
			mock: func(s *mocks.IUserSubscriptionService) {
				s.On("GetAllUserSubscriptions", int64(mockUserContext.ID)).Return(nil, mockError)
			},
			mockUserFromMiddleware: true,
			want: helpers.JsonResponse{
				Code:    http.StatusInternalServerError,
				Message: mockError.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "SUCCESS",
			fields: fields{
				userSubscriptionService: *mocks.NewIUserSubscriptionService(t),
			},
			mock: func(s *mocks.IUserSubscriptionService) {
				s.On("GetAllUserSubscriptions", int64(mockUserContext.ID)).Return([]*models.UserSubscriptions{}, nil)
			},
			mockUserFromMiddleware: true,
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
					UserSubscription: &tt.fields.userSubscriptionService,
				},
			}

			tt.mock(&tt.fields.userSubscriptionService)
			r := helpers.SetUpRouter()
			endpoint := "/user/subscriptions"

			if tt.mockUserFromMiddleware {
				r.GET(endpoint, helpers.MiddlewareMockUser(mockUserContext), h.GetUserSubscriptions)
			} else {
				r.GET(endpoint, h.GetUserSubscriptions)
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

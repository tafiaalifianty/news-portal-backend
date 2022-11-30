package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"
	"final-project-backend/internal/services"
	"final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetAllSubscriptions(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockSubscriptions := []*models.Subscription{}

	type fields struct {
		subscriptionService mocks.ISubscriptionService
	}
	tests := []struct {
		name   string
		fields fields
		mock   func(*mocks.ISubscriptionService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Error from Subscription.GetAll",
			fields: fields{
				subscriptionService: *mocks.NewISubscriptionService(t),
			},
			mock: func(s *mocks.ISubscriptionService) {
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
				subscriptionService: *mocks.NewISubscriptionService(t),
			},
			mock: func(s *mocks.ISubscriptionService) {
				s.On("GetAll").Return(mockSubscriptions, nil)
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
					Subscription: &tt.fields.subscriptionService,
				},
			}

			tt.mock(&tt.fields.subscriptionService)
			r := helpers.SetUpRouter()

			endpoint := "/subscriptions"
			r.GET(endpoint, h.GetAllSubscriptions)

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

func TestHandler_AddUserSubscription(t *testing.T) {
	invalidRequest := &dtos.AddUserSubscriptionRequest{}
	validRequest := &dtos.AddUserSubscriptionRequest{
		UserID:         1,
		SubscriptionID: 1,
	}
	mockTime := time.Now()
	validResponse := &dtos.AddUserSubscriptionResponse{
		UserID:         validRequest.UserID,
		SubscriptionID: validRequest.SubscriptionID,
		DateStarted:    mockTime,
		DateEnded:      mockTime,
	}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)
	mockError := fmt.Errorf("error")

	type fields struct {
		userSubscriptionService *mocks.IUserSubscriptionService
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args
		mock func(*mocks.IUserSubscriptionService)
		want helpers.JsonResponse
	}{
		{
			name: "ERROR | Invalid Request Body",
			fields: fields{
				userSubscriptionService: mocks.NewIUserSubscriptionService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(invalidRequest),
			},
			mock: func(s *mocks.IUserSubscriptionService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from UserSubscription.AddUserSubscription",
			fields: fields{
				userSubscriptionService: mocks.NewIUserSubscriptionService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IUserSubscriptionService) {
				s.On("AddUserSubscription", validRequest.UserID, validRequest.SubscriptionID).Return(nil, mockError)
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
				userSubscriptionService: mocks.NewIUserSubscriptionService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(s *mocks.IUserSubscriptionService) {
				s.On("AddUserSubscription", validRequest.UserID, validRequest.SubscriptionID).Return(&models.UserSubscriptions{
					UserID:         validRequest.UserID,
					SubscriptionID: validRequest.SubscriptionID,
					DateStarted:    mockTime,
					DateEnded:      mockTime,
				}, nil)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
				Data:    mockValidDataInInterface,
				IsError: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				services: &services.Services{
					UserSubscription: tt.fields.userSubscriptionService,
				},
			}

			tt.mock(tt.fields.userSubscriptionService)
			r := helpers.SetUpRouter()

			endpoint := "/subscriptions"
			r.POST(endpoint, h.AddUserSubscription)

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

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"
	"final-project-backend/internal/services"
	"final-project-backend/mocks"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestHandler_Register(t *testing.T) {
	invalidRequest := &dtos.RegisterRequestDTO{}
	validRequest := &dtos.RegisterRequestDTO{
		Fullname: "Fullname",
		Email:    "Email@email.com",
		Password: "Password",
		Address:  "Address",
	}
	validResponse := &dtos.RegisterResponseDTO{
		Email:    validRequest.Email,
		Address:  validRequest.Address,
		Fullname: validRequest.Fullname,
		ID:       1,
	}
	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)
	mockError := fmt.Errorf("error")

	type fields struct {
		authService *mocks.IAuthService
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func(*mocks.IAuthService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Invalid Request Body",
			fields: fields{
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(invalidRequest),
			},
			mock: func(as *mocks.IAuthService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from AuthService: email already exist",
			fields: fields{
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(as *mocks.IAuthService) {
				as.On("Register", &models.User{
					Fullname: validRequest.Fullname,
					Password: validRequest.Password,
					Email:    validRequest.Email,
					Address:  validRequest.Address,
				}, (*string)(nil)).
					Return(nil, nil, nil, &pgconn.PgError{
						Code: errn.UniqueViolation,
					})
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: errn.ErrEmailAlreadyExist.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from AuthService: other errors",
			fields: fields{
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(as *mocks.IAuthService) {
				as.On("Register", &models.User{
					Fullname: validRequest.Fullname,
					Password: validRequest.Password,
					Email:    validRequest.Email,
					Address:  validRequest.Address,
				}, (*string)(nil)).
					Return(nil, nil, nil, mockError)
			},
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
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(as *mocks.IAuthService) {
				as.On("Register", &models.User{
					Fullname: validRequest.Fullname,
					Password: validRequest.Password,
					Email:    validRequest.Email,
					Address:  validRequest.Address,
				}, (*string)(nil)).
					Return(&models.User{
						Fullname: validRequest.Fullname,
						Password: validRequest.Password,
						Email:    validRequest.Email,
						Address:  validRequest.Address,
						ID:       validResponse.ID,
					}, &validResponse.AccessToken, &validResponse.RefreshToken, nil)
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
					Auth: tt.fields.authService,
				},
			}

			tt.mock(tt.fields.authService)
			r := helpers.SetUpRouter()
			endpoint := "/register"
			r.POST(endpoint, h.Register)
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

func TestHandler_Login(t *testing.T) {
	invalidRequest := &dtos.LoginRequestDTO{}
	validRequest := &dtos.LoginRequestDTO{
		Email:    "Email@email.com",
		Password: "Password",
	}
	validResponse := &dtos.LoginResponseDTO{
		ID:    1,
		Email: validRequest.Email,
	}

	mockValidDataInInterface, err := helpers.StructToMap(validResponse)
	require.NoError(t, err)
	mockError := fmt.Errorf("error")

	type fields struct {
		authService *mocks.IAuthService
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mock   func(*mocks.IAuthService)
		want   helpers.JsonResponse
	}{
		{
			name: "ERROR | Invalid Request Body",
			fields: fields{
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(invalidRequest),
			},
			mock: func(as *mocks.IAuthService) {
			},
			want: helpers.JsonResponse{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from AuthService.Login: ErrRecordNotFound",
			fields: fields{
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(as *mocks.IAuthService) {
				as.On("Login", validRequest.Email, validRequest.Password).Return(nil, nil, nil, gorm.ErrRecordNotFound)
			},
			want: helpers.JsonResponse{
				Code:    http.StatusUnauthorized,
				Message: errn.ErrWrongEmailOrPassword.Error(),
				Data:    nil,
				IsError: true,
			},
		},
		{
			name: "ERROR | Error from AuthService.Login: Others",
			fields: fields{
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(as *mocks.IAuthService) {
				as.On("Login", validRequest.Email, validRequest.Password).Return(nil, nil, nil, mockError)
			},
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
				authService: mocks.NewIAuthService(t),
			},
			args: args{
				body: helpers.MakeRequestBody(validRequest),
			},
			mock: func(as *mocks.IAuthService) {
				as.On("Login", validRequest.Email, validRequest.Password).Return(&models.User{
					ID:    validResponse.ID,
					Email: validResponse.Email,
					Role:  models.Member,
				}, &validResponse.AccessToken, &validResponse.RefreshToken, nil)
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
					Auth: tt.fields.authService,
				},
			}

			tt.mock(tt.fields.authService)
			r := helpers.SetUpRouter()
			endpoint := "/login"
			r.POST(endpoint, h.Login)
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

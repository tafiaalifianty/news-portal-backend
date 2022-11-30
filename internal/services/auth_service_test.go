package services

import (
	"fmt"
	"testing"

	"final-project-backend/internal/models"
	mocks "final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewAuthService(t *testing.T) {
	NewAuthService(&AuthServiceConfig{
		userRepository: mocks.NewIUserRepository(t),
		hasher:         mocks.NewHasher(t),
	})
}

func Test_authService_Register(t *testing.T) {
	mockUser := &models.User{
		Fullname: "fullname",
		Email:    "email@email.com",
		Password: "password",
		Address:  "address",
	}
	mockError := fmt.Errorf("error")

	type fields struct {
		userRepository *mocks.IUserRepository
		hasher         *mocks.Hasher
	}
	type args struct {
		user *models.User
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		mock             func(*mocks.IUserRepository, *mocks.Hasher)
		wantUser         *models.User
		wantAccessToken  *string
		wantRefreshToken *string
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "ERROR | Error from Hasher.HashAndSalt",
			fields: fields{
				userRepository: mocks.NewIUserRepository(t),
				hasher:         mocks.NewHasher(t),
			},
			args: args{
				user: mockUser,
			},
			mock: func(ir *mocks.IUserRepository, h *mocks.Hasher) {
				h.On("HashAndSalt", mockUser.Password).Return("", mockError)
			},
			wantUser:    nil,
			wantErr:     true,
			expectedErr: mockError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{
				userRepository: tt.fields.userRepository,
				hasher:         tt.fields.hasher,
			}

			tt.mock(tt.fields.userRepository, tt.fields.hasher)
			gotUser, gotAccessToken, gotRefreshToken, err := s.Register(tt.args.user, nil)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.wantUser, gotUser)
				assert.Equal(t, tt.wantAccessToken, gotAccessToken)
				assert.Equal(t, tt.wantRefreshToken, gotRefreshToken)
			}
		})
	}
}

func Test_authService_Login(t *testing.T) {
	mockUser := &models.User{
		Fullname: "fullname",
		Email:    "email@email.com",
		Password: "password",
		Address:  "address",
	}
	mockError := fmt.Errorf("error")

	type fields struct {
		userRepository *mocks.IUserRepository
		hasher         *mocks.Hasher
	}
	type args struct {
		Email    string
		Password string
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		mock             func(*mocks.IUserRepository, *mocks.Hasher)
		wantUser         *models.User
		wantAccessToken  *string
		wantRefreshToken *string
		wantErr          bool
		expectedErr      error
	}{
		{
			name: "ERROR | Error from UserRepository.GetByEmail",
			fields: fields{
				userRepository: mocks.NewIUserRepository(t),
				hasher:         mocks.NewHasher(t),
			},
			args: args{
				Email:    mockUser.Email,
				Password: mockUser.Password,
			},
			mock: func(ir *mocks.IUserRepository, h *mocks.Hasher) {
				ir.On("GetByEmail", mockUser.Email).Return(&models.User{}, mockError)
			},
			wantUser:        nil,
			wantAccessToken: nil,
			wantErr:         true,
			expectedErr:     mockError,
		},
		{
			name: "ERROR | Error from Hasher.ComparePasswords",
			fields: fields{
				userRepository: mocks.NewIUserRepository(t),
				hasher:         mocks.NewHasher(t),
			},
			args: args{
				Email:    mockUser.Email,
				Password: mockUser.Password,
			},
			mock: func(ir *mocks.IUserRepository, h *mocks.Hasher) {
				ir.On("GetByEmail", mockUser.Email).Return(mockUser, nil)
				h.On("ComparePasswords", mockUser.Password, []byte(mockUser.Password)).Return(false)
			},
			wantUser:        nil,
			wantAccessToken: nil,
			wantErr:         true,
			expectedErr:     gorm.ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{
				userRepository: tt.fields.userRepository,
				hasher:         tt.fields.hasher,
			}

			tt.mock(tt.fields.userRepository, tt.fields.hasher)
			gotUser, gotAccessToken, gotRefreshToken, err := s.Login(tt.args.Email, tt.args.Password)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.wantUser, gotUser)
				assert.Equal(t, tt.wantAccessToken, gotAccessToken)
				assert.Equal(t, tt.wantRefreshToken, gotRefreshToken)
			}
		})
	}
}

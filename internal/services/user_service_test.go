package services

import (
	"fmt"
	"testing"

	"final-project-backend/internal/models"
	mocks "final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	NewUserService(&UserServiceConfig{
		userRepository: mocks.NewIUserRepository(t),
	})
}

func Test_userService_GetByID(t *testing.T) {
	mockUser := &models.User{
		ID:       1,
		Fullname: "fullname",
		Email:    "email@email.com",
		Password: "password",
		Address:  "address",
	}
	mockError := fmt.Errorf("error")
	type fields struct {
		userRepository *mocks.IUserRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IUserRepository)
		want        *models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from UserRepository.GetByID",
			fields: fields{
				userRepository: mocks.NewIUserRepository(t),
			},
			args: args{
				id: mockUser.ID,
			},
			mock: func(ur *mocks.IUserRepository) {
				ur.On("GetByID", mockUser.ID).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				userRepository: mocks.NewIUserRepository(t),
			},
			args: args{
				id: mockUser.ID,
			},
			mock: func(ur *mocks.IUserRepository) {
				ur.On("GetByID", mockUser.ID).Return(mockUser, nil)
			},
			want:        mockUser,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				userRepository: tt.fields.userRepository,
			}

			tt.mock(tt.fields.userRepository)
			got, err := s.GetByID(tt.args.id)

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

package services

import (
	"fmt"
	"testing"

	"final-project-backend/internal/models"
	mocks "final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewSubscriptionService(t *testing.T) {
	NewSubscriptionService(&SubscriptionServiceConfig{
		subscriptionRepository: mocks.NewISubscriptionRepository(t),
	})
}

func Test_subscriptionService_GetAll(t *testing.T) {
	mockSubscriptions := []*models.Subscription{}
	mockError := fmt.Errorf("error")
	type fields struct {
		subscriptionRepository *mocks.ISubscriptionRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.ISubscriptionRepository)
		want        []*models.Subscription
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from SubscriptionHistory.GetAll",
			fields: fields{
				subscriptionRepository: mocks.NewISubscriptionRepository(t),
			},
			mock: func(r *mocks.ISubscriptionRepository) {
				r.On("GetAll").Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				subscriptionRepository: mocks.NewISubscriptionRepository(t),
			},
			mock: func(r *mocks.ISubscriptionRepository) {
				r.On("GetAll").Return(mockSubscriptions, nil)
			},
			want:        mockSubscriptions,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &subscriptionService{
				subscriptionRepository: tt.fields.subscriptionRepository,
			}

			tt.mock(tt.fields.subscriptionRepository)
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

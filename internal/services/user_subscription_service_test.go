package services

import (
	"fmt"
	"testing"

	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/models"
	mocks "final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func Test_userSubscriptionService_AddUserSubscription(t *testing.T) {
	mockUserID := int64(1)
	mockSubscriptionID := int64(1)
	mockError := fmt.Errorf("error")
	type fields struct {
		userSubscriptionRepository *mocks.IUserSubscriptionRepository
		subscriptionRepository     *mocks.ISubscriptionRepository
	}
	type args struct {
		userID         int64
		subscriptionID int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IUserSubscriptionRepository, *mocks.ISubscriptionRepository)
		want        *models.UserSubscriptions
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from subscriptionRepository.GetByID",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:         mockUserID,
				subscriptionID: mockSubscriptionID,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				sr.On("GetByID", mockSubscriptionID).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "ERROR | Error from userSubscriptionRepository.Insert",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:         mockUserID,
				subscriptionID: mockSubscriptionID,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				sr.On("GetByID", mockSubscriptionID).Return(&models.Subscription{ID: mockSubscriptionID}, nil)
				usr.On("Insert", mock.MatchedBy(func(i interface{}) bool {
					userSubscription := i.(*models.UserSubscriptions)
					return userSubscription.UserID == mockUserID &&
						userSubscription.SubscriptionID == mockSubscriptionID
				})).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:         mockUserID,
				subscriptionID: mockSubscriptionID,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				sr.On("GetByID", mockSubscriptionID).Return(&models.Subscription{ID: mockSubscriptionID}, nil)
				usr.On("Insert", mock.MatchedBy(func(i interface{}) bool {
					userSubscription := i.(*models.UserSubscriptions)
					return userSubscription.UserID == mockUserID &&
						userSubscription.SubscriptionID == mockSubscriptionID
				})).Return(&models.UserSubscriptions{}, nil)
			},
			want:        &models.UserSubscriptions{},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userSubscriptionService{
				subscriptionRepository:     tt.fields.subscriptionRepository,
				userSubscriptionRepository: tt.fields.userSubscriptionRepository,
			}

			tt.mock(tt.fields.userSubscriptionRepository, tt.fields.subscriptionRepository)
			got, err := s.AddUserSubscription(tt.args.userID, tt.args.subscriptionID)

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

func Test_userSubscriptionService_ValidateUserQuota(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockUserID := int64(1)
	mockQuotaNeeded := 2
	type fields struct {
		userSubscriptionRepository *mocks.IUserSubscriptionRepository
		subscriptionRepository     *mocks.ISubscriptionRepository
	}
	type args struct {
		userID      int64
		quotaNeeded int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IUserSubscriptionRepository, *mocks.ISubscriptionRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from userSubscriptionRepository.GetOngoingUserSubscriptions",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:      mockUserID,
				quotaNeeded: mockQuotaNeeded,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				usr.On("GetOngoingUserSubscriptions", mockUserID).Return(nil, mockError)
			},
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "ERROR | 0 subscriptions data",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:      mockUserID,
				quotaNeeded: mockQuotaNeeded,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				usr.On("GetOngoingUserSubscriptions", mockUserID).Return([]*models.UserSubscriptions{}, nil)
			},
			wantErr:     true,
			expectedErr: errn.ErrNotEnoughQuota,
		},
		{
			name: "ERROR | not enough quota from available subscriptions",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:      mockUserID,
				quotaNeeded: mockQuotaNeeded,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				usr.On("GetOngoingUserSubscriptions", mockUserID).Return([]*models.UserSubscriptions{
					{
						RemainingQuota: mockQuotaNeeded - 1,
					},
				}, nil)
			},
			wantErr:     true,
			expectedErr: errn.ErrNotEnoughQuota,
		},
		{
			name: "ERROR | Error from DecrementQuota when quotaneeded > remainingquota",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:      mockUserID,
				quotaNeeded: mockQuotaNeeded,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				usr.On("GetOngoingUserSubscriptions", mockUserID).Return([]*models.UserSubscriptions{
					{
						Model: gorm.Model{
							ID: 1,
						},
						RemainingQuota: mockQuotaNeeded - 1,
					},
					{
						Model: gorm.Model{
							ID: 2,
						},
						RemainingQuota: mockQuotaNeeded,
					},
				}, nil)
				usr.On("DecrementQuota", &models.UserSubscriptions{
					Model: gorm.Model{
						ID: 1,
					},
					RemainingQuota: mockQuotaNeeded - 1,
				}, mockQuotaNeeded-1).Return(nil, mockError)
			},
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "ERROR | Error from DecrementQuota when quotaneeded < remainingquota",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:      mockUserID,
				quotaNeeded: mockQuotaNeeded,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				usr.On("GetOngoingUserSubscriptions", mockUserID).Return([]*models.UserSubscriptions{
					{
						Model: gorm.Model{
							ID: 1,
						},
						RemainingQuota: mockQuotaNeeded - 1,
					},
					{
						Model: gorm.Model{
							ID: 2,
						},
						RemainingQuota: mockQuotaNeeded,
					},
				}, nil)
				usr.On("DecrementQuota", &models.UserSubscriptions{
					Model: gorm.Model{
						ID: 1,
					},
					RemainingQuota: mockQuotaNeeded - 1,
				}, mockQuotaNeeded-1).Return(nil, nil)
				usr.On("DecrementQuota", &models.UserSubscriptions{
					Model: gorm.Model{
						ID: 2,
					},
					RemainingQuota: mockQuotaNeeded,
				}, 1).Return(nil, mockError)
			},
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID:      mockUserID,
				quotaNeeded: mockQuotaNeeded,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository, sr *mocks.ISubscriptionRepository) {
				usr.On("GetOngoingUserSubscriptions", mockUserID).Return([]*models.UserSubscriptions{
					{
						Model: gorm.Model{
							ID: 1,
						},
						RemainingQuota: mockQuotaNeeded - 1,
					},
					{
						Model: gorm.Model{
							ID: 2,
						},
						RemainingQuota: mockQuotaNeeded,
					},
				}, nil)
				usr.On("DecrementQuota", &models.UserSubscriptions{
					Model: gorm.Model{
						ID: 1,
					},
					RemainingQuota: mockQuotaNeeded - 1,
				}, mockQuotaNeeded-1).Return(nil, nil)
				usr.On("DecrementQuota", &models.UserSubscriptions{
					Model: gorm.Model{
						ID: 2,
					},
					RemainingQuota: mockQuotaNeeded,
				}, 1).Return(nil, nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userSubscriptionService{
				subscriptionRepository:     tt.fields.subscriptionRepository,
				userSubscriptionRepository: tt.fields.userSubscriptionRepository,
			}

			tt.mock(tt.fields.userSubscriptionRepository, tt.fields.subscriptionRepository)
			err := s.ValidateUserQuota(tt.args.userID, tt.args.quotaNeeded)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func Test_userSubscriptionService_GetAllUserSubscriptions(t *testing.T) {
	mockUserID := int64(1)
	mockError := fmt.Errorf("error")
	type fields struct {
		userSubscriptionRepository *mocks.IUserSubscriptionRepository
		subscriptionRepository     *mocks.ISubscriptionRepository
	}
	type args struct {
		userID int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IUserSubscriptionRepository)
		want        []*models.UserSubscriptions
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from userSubscriptionRepository.GetAllUserSubscriptions",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID: mockUserID,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository) {
				usr.On("GetAllUserSubscriptions", mockUserID).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				subscriptionRepository:     mocks.NewISubscriptionRepository(t),
				userSubscriptionRepository: mocks.NewIUserSubscriptionRepository(t),
			},
			args: args{
				userID: mockUserID,
			},
			mock: func(usr *mocks.IUserSubscriptionRepository) {
				usr.On("GetAllUserSubscriptions", mockUserID).Return([]*models.UserSubscriptions{}, nil)
			},
			want:        []*models.UserSubscriptions{},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userSubscriptionService{
				userSubscriptionRepository: tt.fields.userSubscriptionRepository,
			}

			tt.mock(tt.fields.userSubscriptionRepository)
			got, err := s.GetAllUserSubscriptions(tt.args.userID)

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

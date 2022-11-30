package services

import (
	"fmt"
	"testing"

	"final-project-backend/internal/models"
	"final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewHistoryService(t *testing.T) {
	NewHistoryService(&HistoryServiceConfig{
		historyRepository: mocks.NewIHistoryRepository(t),
	})
}

func Test_historyService_UpdateOrInsert(t *testing.T) {
	mockUpdateError := fmt.Errorf("update_error")
	mockInsertError := fmt.Errorf("insert_error")
	mockHistory := &models.History{}
	type fields struct {
		historyRepository *mocks.IHistoryRepository
	}
	type args struct {
		history *models.History
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IHistoryRepository)
		want        *models.History
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from HistoryRepository.Update",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				history: mockHistory,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("Update", mockHistory).Return(nil, 0, mockUpdateError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockUpdateError,
		},
		{
			name: "ERROR | Error from HistoryRepository.Insert",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				history: mockHistory,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("Update", mockHistory).Return(nil, 0, nil)
				r.On("Insert", mockHistory).Return(nil, mockInsertError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockInsertError,
		},
		{
			name: "SUCCESS | Success when data already exist and update successful",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				history: mockHistory,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("Update", mockHistory).Return(mockHistory, 1, nil)
			},
			want:        nil,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "SUCCESS | Success when data doesn't exist and insert successful",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				history: mockHistory,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("Update", mockHistory).Return(nil, 0, nil)
				r.On("Insert", mockHistory).Return(mockHistory, nil)
			},
			want:        nil,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &historyService{
				historyRepository: tt.fields.historyRepository,
			}

			tt.mock(tt.fields.historyRepository)
			got, err := s.UpdateOrInsert(tt.args.history)

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

func Test_historyService_GetByUserID(t *testing.T) {
	mockUserID := int64(1)
	mockError := fmt.Errorf("error")

	mockHistories := []*models.History{}
	type fields struct {
		historyRepository *mocks.IHistoryRepository
	}
	type args struct {
		userID int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IHistoryRepository)
		want        []*models.History
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from GetByUserID.Update",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				userID: mockUserID,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("GetByUserID", mockUserID).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				userID: mockUserID,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("GetByUserID", mockUserID).Return(mockHistories, nil)
			},
			want:        mockHistories,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &historyService{
				historyRepository: tt.fields.historyRepository,
			}

			tt.mock(tt.fields.historyRepository)
			got, err := s.GetByUserID(tt.args.userID)

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

func Test_historyService_GetByUserAndPostID(t *testing.T) {
	mockUserID := int64(1)
	mockPostID := int64(1)
	mockError := fmt.Errorf("error")

	mockHistories := &models.History{}
	type fields struct {
		historyRepository *mocks.IHistoryRepository
	}
	type args struct {
		userID int64
		postID int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IHistoryRepository)
		want        *models.History
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from GetByUserAndPostID.Update",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				userID: mockUserID,
				postID: mockPostID,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("GetByUserAndPostID", mockUserID, mockPostID).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				historyRepository: mocks.NewIHistoryRepository(t),
			},
			args: args{
				userID: mockUserID,
				postID: mockPostID,
			},
			mock: func(r *mocks.IHistoryRepository) {
				r.On("GetByUserAndPostID", mockUserID, mockPostID).Return(mockHistories, nil)
			},
			want:        mockHistories,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &historyService{
				historyRepository: tt.fields.historyRepository,
			}

			tt.mock(tt.fields.historyRepository)
			got, err := s.GetByUserAndPostID(tt.args.userID, tt.args.postID)

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

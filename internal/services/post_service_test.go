package services

import (
	"fmt"
	"testing"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/models"
	mocks "final-project-backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewPostService(t *testing.T) {
	NewPostService(&PostServiceConfig{
		postRepository: mocks.NewIPostRepository(t),
	})
}

func Test_postService_Create(t *testing.T) {
	mockPostParams := &models.Post{}
	mockPostReturn := &models.Post{ID: 1}

	mockError := fmt.Errorf("error")
	type fields struct {
		postRepository *mocks.IPostRepository
	}
	type args struct {
		post *models.Post
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IPostRepository)
		want        *models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from PostRepository.Insert",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			args: args{
				post: mockPostParams,
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("Insert", mockPostParams).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			args: args{
				post: mockPostParams,
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("Insert", mockPostParams).Return(mockPostReturn, nil)
			},
			want:        mockPostReturn,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &postService{
				postRepository: tt.fields.postRepository,
			}

			tt.mock(tt.fields.postRepository)
			got, err := s.Create(tt.args.post)

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

func Test_postService_GetByID(t *testing.T) {
	mockPost := &models.Post{
		ID: 1,
	}
	mockError := fmt.Errorf("error")
	type fields struct {
		postRepository *mocks.IPostRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IPostRepository)
		want        *models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from PostRepository.GetByID",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			args: args{
				id: mockPost.ID,
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetByID", mockPost.ID).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			args: args{
				id: mockPost.ID,
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetByID", mockPost.ID).Return(mockPost, nil)
			},
			want:        mockPost,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &postService{
				postRepository: tt.fields.postRepository,
			}

			tt.mock(tt.fields.postRepository)
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

func Test_postService_GetAll(t *testing.T) {
	mockPosts := []*models.Post{
		{ID: 1},
	}
	mockError := fmt.Errorf("error")
	mockQuery := &dtos.PostsRequestQuery{}
	type fields struct {
		postRepository *mocks.IPostRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.IPostRepository)
		want        []*models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from PostRepository.GetAll",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetAll", mockQuery).Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetAll", mockQuery).Return(mockPosts, nil)
			},
			want:        mockPosts,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &postService{
				postRepository: tt.fields.postRepository,
			}

			tt.mock(tt.fields.postRepository)
			got, err := s.GetAll(mockQuery)

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

func Test_postService_Delete(t *testing.T) {
	mockPost := &models.Post{
		ID: 1,
	}
	mockError := fmt.Errorf("error")
	type fields struct {
		postRepository *mocks.IPostRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func(*mocks.IPostRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from PostRepository.Delete",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			args: args{
				id: mockPost.ID,
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("Delete", mockPost.ID).Return(mockError)
			},
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			args: args{
				id: mockPost.ID,
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("Delete", mockPost.ID).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &postService{
				postRepository: tt.fields.postRepository,
			}

			tt.mock(tt.fields.postRepository)
			err := s.Delete(tt.args.id)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func Test_postService_CountByQuery(t *testing.T) {
	mockError := fmt.Errorf("error")
	mockQuery := &dtos.PostsRequestQuery{}
	type fields struct {
		postRepository *mocks.IPostRepository
	}
	tests := []struct {
		name           string
		fields         fields
		mock           func(*mocks.IPostRepository)
		wantTotalPosts int64
		wantTotalPages int64
		wantErr        bool
		expectedErr    error
	}{
		{
			name: "ERROR | Error from PostRepository.CountByQuery",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("CountByQuery", mockQuery).Return(int64(0), mockError)
			},
			wantTotalPosts: 0,
			wantTotalPages: 0,
			wantErr:        true,
			expectedErr:    mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("CountByQuery", mockQuery).Return(int64(1), nil)
			},
			wantTotalPosts: 1,
			wantTotalPages: 1,
			wantErr:        false,
			expectedErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &postService{
				postRepository: tt.fields.postRepository,
			}

			tt.mock(tt.fields.postRepository)
			gotTotalPosts, gotTotalPages, err := s.CountByQuery(mockQuery)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.wantTotalPosts, gotTotalPosts)
				assert.Equal(t, tt.wantTotalPages, gotTotalPages)
			}
		})
	}
}

func Test_postService_GetTypes(t *testing.T) {
	mockError := fmt.Errorf("error")
	type fields struct {
		postRepository *mocks.IPostRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.IPostRepository)
		want        []*models.PostType
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from PostRepository.GetAll",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetTypes").Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetTypes").Return([]*models.PostType{}, nil)
			},
			want:        []*models.PostType{},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &postService{
				postRepository: tt.fields.postRepository,
			}

			tt.mock(tt.fields.postRepository)
			got, err := s.GetTypes()

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

func Test_postService_GetCategories(t *testing.T) {
	mockError := fmt.Errorf("error")
	type fields struct {
		postRepository *mocks.IPostRepository
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func(*mocks.IPostRepository)
		want        []*models.Category
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | Error from PostRepository.GetAll",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetCategories").Return(nil, mockError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: mockError,
		},
		{
			name: "SUCCESS",
			fields: fields{
				postRepository: mocks.NewIPostRepository(t),
			},
			mock: func(r *mocks.IPostRepository) {
				r.On("GetCategories").Return([]*models.Category{}, nil)
			},
			want:        []*models.Category{},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &postService{
				postRepository: tt.fields.postRepository,
			}

			tt.mock(tt.fields.postRepository)
			got, err := s.GetCategories()

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

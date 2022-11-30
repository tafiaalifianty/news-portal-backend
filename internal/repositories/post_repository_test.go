package repositories

import (
	"database/sql"
	"testing"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewPostRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	NewPostRepository(&PostRepositoryConfig{
		db: DB,
	})
}

func Test_postRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	mockPost := &models.Post{}
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		post *models.Post
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func()
		want        *models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR",
			fields: fields{
				db: DB,
			},
			args: args{
				post: mockPost,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT (.+)").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "SUCCESS",
			fields: fields{
				db: DB,
			},
			args: args{
				post: mockPost,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT (.+)").
					WillReturnRows(sqlmock.NewRows(([]string{"id"})).AddRow(1))

				mock.ExpectCommit()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:        mockPost,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{
				db: tt.fields.db,
			}

			tt.mock()
			got, err := r.Insert(tt.args.post)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_postRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	mockPost := &models.Post{}
	mockPosts := []*models.Post{mockPost}
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name        string
		fields      fields
		mock        func()
		want        []*models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR",
			fields: fields{
				db: DB,
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+)").
					WillReturnError(sql.ErrNoRows)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "SUCCESS",
			fields: fields{
				db: DB,
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+)").
					WillReturnRows(sqlmock.NewRows(([]string{"id"})).AddRow(mockPost.ID))

				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(mockPost.ID, 1))
			},
			want:        mockPosts,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{
				db: tt.fields.db,
			}

			tt.mock()
			got, err := r.GetAll(&dtos.PostsRequestQuery{})

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_postRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	mockID := 1
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ID int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func()
		want        *models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR",
			fields: fields{
				db: DB,
			},
			args: args{
				ID: int64(mockID),
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+)").
					WithArgs(mockID).
					WillReturnError(sql.ErrNoRows)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "SUCCESS",
			fields: fields{
				db: DB,
			},
			args: args{
				ID: int64(mockID),
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+)").
					WithArgs(mockID).
					WillReturnRows(sqlmock.NewRows(([]string{"id"})).AddRow(mockID))

				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(int64(mockID), 1))
			},
			want:        &models.Post{ID: int64(mockID)},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{
				db: tt.fields.db,
			}

			tt.mock()
			got, err := r.GetByID(tt.args.ID)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

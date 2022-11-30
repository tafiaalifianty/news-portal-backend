package repositories

import (
	"database/sql"
	"testing"

	"final-project-backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	NewUserRepository(&UserRepositoryConfig{
		db: DB,
	})
}

func Test_userRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	mockUser := &models.User{}
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user *models.User
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func()
		want        *models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR",
			fields: fields{
				db: DB,
			},
			args: args{
				user: mockUser,
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
				user: mockUser,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT (.+)").
					WillReturnRows(sqlmock.NewRows(([]string{"id"})).AddRow(1))

				mock.ExpectCommit()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:        mockUser,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: tt.fields.db,
			}

			tt.mock()
			got, err := r.Insert(tt.args.user)

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

func Test_userRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	mockEmail := "test@email.com"
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func()
		want        *models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | No user with email found",
			fields: fields{
				db: DB,
			},
			args: args{
				email: mockEmail,
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+)").
					WithArgs(mockEmail).
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
				email: mockEmail,
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+)").
					WithArgs(mockEmail).
					WillReturnRows(sqlmock.NewRows(([]string{"email"})).AddRow(mockEmail))

				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:        &models.User{Email: mockEmail},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: tt.fields.db,
			}

			tt.mock()
			got, err := r.GetByEmail(tt.args.email)

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

func Test_userRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	var mockID int64 = 1
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mock        func()
		want        *models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "ERROR | No user with email found",
			fields: fields{
				db: DB,
			},
			args: args{
				id: mockID,
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
				id: mockID,
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+)").
					WithArgs(mockID).
					WillReturnRows(sqlmock.NewRows(([]string{"id"})).AddRow(mockID))

				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:        &models.User{ID: mockID},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{
				db: tt.fields.db,
			}

			tt.mock()
			got, err := r.GetByID(tt.args.id)

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

package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)

	dia := postgres.New(postgres.Config{
		DriverName: "postgres",
		Conn:       db,
	})
	DB, err := gorm.Open(dia)
	assert.NoError(t, err)

	New(DB)
}

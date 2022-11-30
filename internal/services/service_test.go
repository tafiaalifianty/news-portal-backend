package services

import (
	"testing"

	"final-project-backend/internal/repositories"
)

func TestNew(t *testing.T) {
	New(&repositories.Repositories{})
}

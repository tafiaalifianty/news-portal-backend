package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
)

type IUserGiftRepository interface {
	Insert(userGift *models.UserGift) (*models.UserGift, error)
}

type userGiftRepository struct {
	db *gorm.DB
}

type UserGiftRepositoryConfig struct {
	db *gorm.DB
}

func NewUserGiftRepository(c *UserGiftRepositoryConfig) IUserGiftRepository {
	return &userGiftRepository{
		db: c.db,
	}
}

func (r *userGiftRepository) Insert(userGift *models.UserGift) (*models.UserGift, error) {
	result := r.db.Create(userGift)
	if result.Error != nil {
		return nil, result.Error
	}

	return userGift, nil
}

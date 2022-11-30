package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
)

type IUserSpendingRepository interface {
	Insert(userSpending *models.UserSpending) (*models.UserSpending, error)
	GetByPrimaryKeys(userID int64, month int, year int) (*models.UserSpending, error)
	Update(userSpending *models.UserSpending) (*models.UserSpending, int, error)
}

type userSpendingRepository struct {
	db *gorm.DB
}

type UserSpendingRepositoryConfig struct {
	db *gorm.DB
}

func NewUserSpendingRepository(c *UserSpendingRepositoryConfig) IUserSpendingRepository {
	return &userSpendingRepository{
		db: c.db,
	}
}

func (r *userSpendingRepository) Insert(userSpending *models.UserSpending) (*models.UserSpending, error) {
	result := r.db.Create(userSpending)
	if result.Error != nil {
		return nil, result.Error
	}

	return userSpending, nil
}

func (r *userSpendingRepository) GetByPrimaryKeys(userID int64, month int, year int) (*models.UserSpending, error) {
	var userSpending *models.UserSpending

	result := r.db.
		Where("user_id = ? AND month = ? AND year = ?", userID, month, year).
		First(&userSpending)

	if result.Error != nil {
		return nil, result.Error
	}

	return userSpending, nil
}

func (r *userSpendingRepository) Update(userSpending *models.UserSpending) (*models.UserSpending, int, error) {
	result := r.db.Model(&userSpending).Updates(&userSpending)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return userSpending, int(result.RowsAffected), nil
}

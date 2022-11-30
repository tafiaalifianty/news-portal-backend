package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserTokenRepository interface {
	Insert(userToken *models.UserToken) (*models.UserToken, error)
	GetByUserID(userID int64) (*models.UserToken, error)
	Update(token *models.UserToken) (*models.UserToken, int, error)
}

type userTokenRepository struct {
	db *gorm.DB
}

type UserTokenRepositoryConfig struct {
	db *gorm.DB
}

func NewUserTokenRepository(c *UserTokenRepositoryConfig) IUserTokenRepository {
	return &userTokenRepository{db: c.db}
}

func (r *userTokenRepository) Insert(userToken *models.UserToken) (*models.UserToken, error) {
	result := r.db.Create(&userToken)
	if result.Error != nil {
		return nil, result.Error
	}

	return userToken, nil
}

func (r *userTokenRepository) GetByUserID(userID int64) (*models.UserToken, error) {
	var userToken *models.UserToken

	result := r.db.
		Where("user_id = ?", userID).
		First(&userToken)

	if result.Error != nil {
		return nil, result.Error
	}

	return userToken, nil
}

func (r *userTokenRepository) Update(token *models.UserToken) (*models.UserToken, int, error) {
	result := r.db.Model(&token).Clauses(clause.Returning{}).Updates(&token)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return token, int(result.RowsAffected), nil
}

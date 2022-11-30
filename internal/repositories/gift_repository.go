package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IGiftRepository interface {
	GetAll() ([]*models.Gift, error)
	GetAllUserGifts() ([]*models.UserGift, error)
	GetUserGift(userGift *models.UserGift) (*models.UserGift, error)
	GetUserGiftsByUserID(userID int64) ([]*models.UserGift, error)
	UpdateStock(gift *models.Gift) (*models.Gift, int, error)
	UpdateUserGift(userGift *models.UserGift) (*models.UserGift, int, error)
}

type giftRepository struct {
	db *gorm.DB
}

type GiftRepositoryConfig struct {
	db *gorm.DB
}

func NewGiftRepository(c *GiftRepositoryConfig) IGiftRepository {
	return &giftRepository{
		db: c.db,
	}
}

func (r *giftRepository) GetAll() ([]*models.Gift, error) {
	var gifts []*models.Gift
	result := r.db.Select("id", "name", "minimum_spending", "stock", "updated_at").Order("id asc").Find(&gifts)

	if result.Error != nil {
		return nil, result.Error
	}

	return gifts, nil
}

func (r *giftRepository) GetAllUserGifts() ([]*models.UserGift, error) {
	var userGifts []*models.UserGift
	result := r.db.Joins("Gift").Joins("User").Find(&userGifts)

	if result.Error != nil {
		return nil, result.Error
	}

	return userGifts, nil
}

func (r *giftRepository) GetUserGift(userGift *models.UserGift) (*models.UserGift, error) {

	result := r.db.
		First(&userGift)

	if result.Error != nil {
		return nil, result.Error
	}

	return userGift, nil
}

func (r *giftRepository) GetUserGiftsByUserID(userID int64) ([]*models.UserGift, error) {
	var userGifts []*models.UserGift

	result := r.db.
		Where("user_id = ?", userID).
		Joins("Gift").
		Find(&userGifts)

	if result.Error != nil {
		return nil, result.Error
	}

	return userGifts, nil
}

func (r *giftRepository) UpdateStock(gift *models.Gift) (*models.Gift, int, error) {
	result := r.db.Model(&gift).Select("stock").Clauses(clause.Returning{}).Updates(&gift)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return gift, int(result.RowsAffected), nil
}

func (r *giftRepository) UpdateUserGift(userGift *models.UserGift) (*models.UserGift, int, error) {
	result := r.db.Model(&userGift).Updates(&userGift)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return userGift, int(result.RowsAffected), nil
}

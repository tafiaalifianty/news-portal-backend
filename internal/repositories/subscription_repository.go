package repositories

import (
	"final-project-backend/internal/models"

	"gorm.io/gorm"
)

type ISubscriptionRepository interface {
	GetAll() ([]*models.Subscription, error)
	GetByID(id int64) (*models.Subscription, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

type SubscriptionRepositoryConfig struct {
	db *gorm.DB
}

func NewSubscriptionRepository(c *SubscriptionRepositoryConfig) ISubscriptionRepository {
	return &subscriptionRepository{
		db: c.db,
	}
}

func (r *subscriptionRepository) GetAll() ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription
	result := r.db.Find(&subscriptions)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}

func (r *subscriptionRepository) GetByID(id int64) (*models.Subscription, error) {
	var subscription *models.Subscription

	result := r.db.
		Where("subscriptions.id = ?", id).
		First(&subscription)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}

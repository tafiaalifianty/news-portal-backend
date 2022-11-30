package repositories

import (
	"time"

	"final-project-backend/internal/models"

	"gorm.io/gorm"
)

type IUserSubscriptionRepository interface {
	Insert(userSubscription *models.UserSubscriptions) (*models.UserSubscriptions, error)
	GetOngoingUserSubscriptions(userID int64) ([]*models.UserSubscriptions, error)
	GetAllUserSubscriptions(userID int64) ([]*models.UserSubscriptions, error)
	Update(userSubscription *models.UserSubscriptions) (*models.UserSubscriptions, error)
	DecrementQuota(userSubscription *models.UserSubscriptions, quota int) (*models.UserSubscriptions, error)
}

type userSubscriptionRepository struct {
	db *gorm.DB
}

type UserSubscriptionRepositoryConfig struct {
	db *gorm.DB
}

func NewUserSubscriptionRepository(c *UserSubscriptionRepositoryConfig) IUserSubscriptionRepository {
	return &userSubscriptionRepository{db: c.db}
}

func (r *userSubscriptionRepository) Insert(userSubscription *models.UserSubscriptions) (*models.UserSubscriptions, error) {
	result := r.db.Create(&userSubscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return userSubscription, nil
}

func (r *userSubscriptionRepository) GetAllUserSubscriptions(userID int64) ([]*models.UserSubscriptions, error) {
	var userSubscriptions []*models.UserSubscriptions

	result := r.db.
		Where("user_id = ?", userID).
		Order("date_ended asc, remaining_quota desc").
		Find(&userSubscriptions)

	if result.Error != nil {
		return nil, result.Error
	}

	return userSubscriptions, nil
}

func (r *userSubscriptionRepository) GetOngoingUserSubscriptions(userID int64) ([]*models.UserSubscriptions, error) {
	var userSubscriptions []*models.UserSubscriptions

	result := r.db.
		Where("user_id = ?", userID).
		Where("date_ended > ?", time.Now()).
		Where("remaining_quota > ?", 0).
		Order("date_ended asc, remaining_quota desc").
		Find(&userSubscriptions)

	if result.Error != nil {
		return nil, result.Error
	}

	return userSubscriptions, nil
}

func (r *userSubscriptionRepository) Update(userSubscription *models.UserSubscriptions) (*models.UserSubscriptions, error) {
	result := r.db.Model(userSubscription).Updates(userSubscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return userSubscription, nil
}

func (r *userSubscriptionRepository) DecrementQuota(userSubscription *models.UserSubscriptions, quota int) (*models.UserSubscriptions, error) {
	result := r.db.Model(userSubscription).Update("remaining_quota", userSubscription.RemainingQuota-quota)
	if result.Error != nil {
		return nil, result.Error
	}

	return userSubscription, nil
}

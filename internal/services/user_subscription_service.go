package services

import (
	"time"

	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"
)

type IUserSubscriptionService interface {
	AddUserSubscription(userID int64, subscriptionID int64) (*models.UserSubscriptions, error)
	GetAllUserSubscriptions(userID int64) ([]*models.UserSubscriptions, error)
	ValidateUserQuota(userID int64, quotaNeeded int) error
}

type userSubscriptionService struct {
	userSubscriptionRepository repositories.IUserSubscriptionRepository
	subscriptionRepository     repositories.ISubscriptionRepository
}

type UserSubscriptionServiceConfig struct {
	userSubscriptionRepository repositories.IUserSubscriptionRepository
	subscriptionRepository     repositories.ISubscriptionRepository
}

func NewUserSubscriptionService(c *UserSubscriptionServiceConfig) IUserSubscriptionService {
	return &userSubscriptionService{
		userSubscriptionRepository: c.userSubscriptionRepository,
		subscriptionRepository:     c.subscriptionRepository,
	}
}

func (s *userSubscriptionService) AddUserSubscription(userID int64, subscriptionID int64) (*models.UserSubscriptions, error) {
	subscription, err := s.subscriptionRepository.GetByID(subscriptionID)

	if err != nil {
		return nil, err
	}

	userSubscription := &models.UserSubscriptions{
		UserID:         userID,
		SubscriptionID: subscriptionID,
		RemainingQuota: subscription.Quota,
		DateStarted:    time.Now(),
		DateEnded:      time.Now().AddDate(0, 1, 0),
	}

	createdUserSubscription, err := s.userSubscriptionRepository.Insert(userSubscription)
	if err != nil {
		return nil, err
	}

	return createdUserSubscription, nil
}

func (s *userSubscriptionService) GetAllUserSubscriptions(userID int64) ([]*models.UserSubscriptions, error) {
	userSubscriptions, err := s.userSubscriptionRepository.GetAllUserSubscriptions(userID)
	if err != nil {
		return nil, err
	}

	return userSubscriptions, nil
}

func (s *userSubscriptionService) ValidateUserQuota(userID int64, quotaNeeded int) error {
	userSubscriptions, err := s.userSubscriptionRepository.GetOngoingUserSubscriptions(userID)

	if err != nil {
		return err
	}

	totalQuota := 0
	for _, subscription := range userSubscriptions {
		totalQuota += subscription.RemainingQuota
	}

	if totalQuota < quotaNeeded {
		return errn.ErrNotEnoughQuota
	}

	for _, subscription := range userSubscriptions {
		if quotaNeeded > subscription.RemainingQuota {
			quotaNeeded -= subscription.RemainingQuota
			_, err := s.userSubscriptionRepository.DecrementQuota(subscription, subscription.RemainingQuota)

			if err != nil {
				return err
			}
		} else {
			_, err := s.userSubscriptionRepository.DecrementQuota(subscription, quotaNeeded)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

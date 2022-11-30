package services

import (
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"
)

type ISubscriptionService interface {
	GetAll() ([]*models.Subscription, error)
}

type subscriptionService struct {
	subscriptionRepository repositories.ISubscriptionRepository
}

type SubscriptionServiceConfig struct {
	subscriptionRepository repositories.ISubscriptionRepository
}

func NewSubscriptionService(c *SubscriptionServiceConfig) ISubscriptionService {
	return &subscriptionService{subscriptionRepository: c.subscriptionRepository}
}

func (s *subscriptionService) GetAll() ([]*models.Subscription, error) {
	subscriptions, err := s.subscriptionRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

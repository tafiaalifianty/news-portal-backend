package dtos

import (
	"time"

	"final-project-backend/internal/models"
)

type GetAllSubscriptionsResponse = []*SubscriptionResponse

type SubscriptionResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Quota int    `json:"quota"`
}

type AddUserSubscriptionRequest struct {
	UserID         int64 `json:"user_id" binding:"required"`
	SubscriptionID int64 `json:"subscription_id" binding:"required"`
}

type AddUserSubscriptionResponse struct {
	UserID         int64     `json:"user_id"`
	SubscriptionID int64     `json:"subscription_id"`
	RemainingQuota int       `json:"remaining_quota"`
	DateStarted    time.Time `json:"date_started"`
	DateEnded      time.Time `json:"date_ended"`
}

func FormatSubscription(s *models.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		ID:    s.ID,
		Name:  s.Name,
		Price: s.Price,
		Quota: s.Quota,
	}
}

func FormatSubscriptions(subscriptions []*models.Subscription) []*SubscriptionResponse {
	formattedSubscriptions := []*SubscriptionResponse{}
	for _, subscription := range subscriptions {
		formattedSubscription := FormatSubscription(subscription)
		formattedSubscriptions = append(formattedSubscriptions, formattedSubscription)
	}
	return formattedSubscriptions
}

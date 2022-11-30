package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSubscriptions struct {
	gorm.Model
	UserID         int64        `json:"user_id"`
	SubscriptionID int64        `json:"subscription_id"`
	Subscription   Subscription `json:"subscription" gorm:"foreignKey:subscription_id"`
	RemainingQuota int          `json:"remaining_quota"`
	DateStarted    time.Time    `json:"date_started"`
	DateEnded      time.Time    `json:"date_ended"`
}

func (UserSubscriptions) BeforeCreate(db *gorm.DB) error {
	return nil
}

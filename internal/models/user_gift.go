package models

import (
	"final-project-backend/internal/constants"
)

type UserGift struct {
	UserID int64                `json:"user_id" gorm:"primaryKey"`
	User   User                 `json:"user"`
	GiftID int64                `json:"gift_id" gorm:"primaryKey"`
	Gift   Gift                 `json:"gift"`
	Month  int                  `json:"month" gorm:"primaryKey"`
	Year   int                  `json:"year" gorm:"primaryKey"`
	Status constants.GiftStatus `json:"status"`
}

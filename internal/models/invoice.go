package models

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Code           string        `json:"string" gorm:"unique"`
	UserID         int64         `json:"user_id"`
	User           User          `json:"user" gorm:"foreignKey:user_id"`
	Status         InvoiceStatus `json:"status"`
	Total          int           `json:"total"`
	OriginalPrice  int           `json:"original_price"`
	PaidAt         *time.Time    `json:"paid_at" gorm:"default:null"`
	SubscriptionID int64         `json:"subscription_id"`
	Subscription   Subscription  `json:"subscription" gorm:"foreignKey:subscription_id"`
	VoucherCode    string        `json:"voucher_code"`
}

type InvoiceStatus int

const (
	WAITING InvoiceStatus = iota + 1
	PROCESSED
	COMPLETED
	REJECTED
)

func (e InvoiceStatus) String() string {
	switch e {
	case WAITING:
		return "WAITING"
	case PROCESSED:
		return "PROCESSED"
	case COMPLETED:
		return "COMPLETED"
	case REJECTED:
		return "REJECTED"
	default:
		return ""
	}
}

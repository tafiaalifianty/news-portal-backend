package models

import (
	"time"

	"final-project-backend/internal/constants"
)

type UserVoucher struct {
	ID                 int64                   `json:"id" gorm:"primaryKey"`
	Code               string                  `json:"code" gorm:"unique"`
	UserID             int64                   `json:"user_id"`
	ReceivedFromUserID int64                   `json:"received_from_user_id"`
	ReceivedFromUser   User                    `json:"user" gorm:"foreignKey:received_from_user_id"`
	VoucherID          int64                   `json:"voucher_id"`
	Voucher            Voucher                 `json:"voucher" gorm:"foreignKey:voucher_id"`
	DateReceived       time.Time               `json:"date_received"`
	ValidUntil         time.Time               `json:"valid_until"`
	Status             constants.VoucherStatus `json:"status"`
}

package dtos

import (
	"time"

	"final-project-backend/internal/constants"
	"final-project-backend/internal/models"

	"gorm.io/gorm"
)

type GetUserProfileResponseDTO struct {
	ID           int64       `json:"id"`
	Email        string      `json:"email"`
	Role         models.Role `json:"role"`
	Fullname     string      `json:"fullname"`
	Address      string      `json:"address"`
	ReferralCode string      `json:"referral_code"`
}

type HistoryResponseDTO struct {
	PostID       int64                `json:"post_id"`
	LastAccessed time.Time            `json:"last_accessed" gorm:"autoUpdateTime:true"`
	Post         *PostResponseCompact `json:"post,omitempty"`
	IsLiked      bool                 `json:"is_liked"`
	IsShared     bool                 `json:"is_shared"`
	DeletedAt    gorm.DeletedAt
	CreatedAt    time.Time
}

type GetUserHistoriesDTO = []*HistoryResponseDTO

type UserSubscriptionDTO struct {
	ID             int64     `json:"id"`
	SubscriptionID int64     `json:"subscription_id"`
	RemainingQuota int       `json:"remaining_quota"`
	DateStarted    time.Time `json:"date_started"`
	DateEnded      time.Time `json:"date_ended"`
}

type GetUserSubscriptionsDTO = []*UserSubscriptionDTO

type GetUserGiftsResponse = []*UserGiftResponse

type UserReferralResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type UserVoucherResponse struct {
	ID                    int64          `json:"id" gorm:"primaryKey"`
	Code                  string         `json:"code" gorm:"unique"`
	ReceivedFromUserID    int64          `json:"received_from_user_id"`
	ReceivedFromUserEmail string         `json:"received_from_user_email"`
	VoucherID             int64          `json:"voucher_id"`
	Voucher               models.Voucher `json:"voucher" gorm:"foreignKey:voucher_id"`
	DateReceived          time.Time      `json:"date_received"`
	ValidUntil            time.Time      `json:"valid_until"`
	IsValid               bool           `json:"is_valid"`
	Status                string         `json:"status"`
}

type GetUsersReferralsResponse = []*models.UserSpendingTotalAggregates

type GetUserVouchersResponse = []*UserVoucherResponse

type UpdateUserRequest struct {
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
}

type UpdateUserResponse struct {
	ID       int64  `gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique"`
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
}

func FormatUpdattedUser(user *models.User) *UpdateUserResponse {
	return &UpdateUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
		Address:  user.Address,
	}
}

func FormatUserVoucher(userVoucher *models.UserVoucher) *UserVoucherResponse {
	return &UserVoucherResponse{
		ID:                    userVoucher.ID,
		Code:                  userVoucher.Code,
		ReceivedFromUserID:    userVoucher.ReceivedFromUserID,
		ReceivedFromUserEmail: userVoucher.ReceivedFromUser.Email,
		VoucherID:             userVoucher.VoucherID,
		Voucher:               userVoucher.Voucher,
		DateReceived:          userVoucher.DateReceived,
		ValidUntil:            userVoucher.ValidUntil,
		IsValid:               time.Now().Before(userVoucher.ValidUntil) && userVoucher.Status == constants.AVAILABLE,
		Status:                userVoucher.Status.String(),
	}
}

func FormatUserVouchers(userVouchers []*models.UserVoucher) []*UserVoucherResponse {
	formattedUserVouchers := []*UserVoucherResponse{}
	for _, userVoucher := range userVouchers {
		formatted := FormatUserVoucher(userVoucher)
		formattedUserVouchers = append(formattedUserVouchers, formatted)
	}

	return formattedUserVouchers
}

func FormatUserReferral(user *models.User) *UserReferralResponse {
	return &UserReferralResponse{
		ID:       user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
	}
}

func FormatUserReferrals(users []*models.User) []*UserReferralResponse {
	formattedUserReferrals := []*UserReferralResponse{}
	for _, user := range users {
		formatted := FormatUserReferral(user)
		formattedUserReferrals = append(formattedUserReferrals, formatted)
	}

	return formattedUserReferrals
}

func FormatHistory(history *models.History) *HistoryResponseDTO {
	return &HistoryResponseDTO{
		PostID:       history.PostID,
		LastAccessed: history.LastAccessed,
		IsLiked:      history.IsLiked,
		IsShared:     history.IsShared,
		Post:         FormatPostCompact(&history.Post),
	}
}

func FormatHistories(histories []*models.History) []*HistoryResponseDTO {
	formattedHistories := []*HistoryResponseDTO{}
	for _, history := range histories {
		formattedHistory := FormatHistory(history)
		formattedHistories = append(formattedHistories, formattedHistory)
	}
	return formattedHistories
}

func FormatUserSubscription(userSubscription *models.UserSubscriptions) *UserSubscriptionDTO {
	return &UserSubscriptionDTO{
		ID:             int64(userSubscription.ID),
		SubscriptionID: userSubscription.SubscriptionID,
		RemainingQuota: userSubscription.RemainingQuota,
		DateStarted:    userSubscription.DateStarted,
		DateEnded:      userSubscription.DateEnded,
	}
}

func FormatUserSubscriptions(userSubscriptions []*models.UserSubscriptions) []*UserSubscriptionDTO {
	formattedUserSubscriptions := []*UserSubscriptionDTO{}
	for _, userSubscription := range userSubscriptions {
		formattedUserSubscription := FormatUserSubscription(userSubscription)
		formattedUserSubscriptions = append(formattedUserSubscriptions, formattedUserSubscription)
	}
	return formattedUserSubscriptions
}

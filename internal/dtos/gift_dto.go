package dtos

import (
	"time"

	"final-project-backend/internal/models"
)

type GiftResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Stock     int        `json:"stock"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type GetAllGiftsResponse = []*GiftResponse

type UserGiftResponse struct {
	UserID    int64  `json:"user_id" gorm:"primaryKey"`
	UserEmail string `json:"user_email,omitempty"`
	GiftID    int64  `json:"gift_id" gorm:"primaryKey"`
	GiftName  string `json:"gift_name,omitempty"`
	Month     int    `json:"month" gorm:"primaryKey"`
	Year      int    `json:"year" gorm:"primaryKey"`
	Status    string `json:"status"`
}

type GetAllUserGiftsResponse = []*UserGiftResponse

type UpdateUserGiftStatusRequest struct {
	IsSuccess *bool `json:"is_success" binding:"required"`
}

type UpdateUserGiftStatusResponse = *UserGiftResponse

type UpdateGiftStockRequest struct {
	Stock int `json:"stock" binding:"required"`
}

type UpdateGiftStockResponse = *GiftResponse

func FormatGift(gift *models.Gift) *GiftResponse {
	return &GiftResponse{
		ID:        int64(gift.ID),
		Name:      gift.Name,
		Stock:     gift.Stock,
		UpdatedAt: &gift.UpdatedAt,
	}
}

func FormatGifts(gifts []*models.Gift) []*GiftResponse {
	formattedGifts := []*GiftResponse{}
	for _, gift := range gifts {
		formattedGift := FormatGift(gift)
		formattedGifts = append(formattedGifts, formattedGift)
	}
	return formattedGifts
}

func FormatUserGift(userGift *models.UserGift) *UserGiftResponse {
	return &UserGiftResponse{
		UserID:    userGift.UserID,
		UserEmail: userGift.User.Email,
		GiftID:    userGift.GiftID,
		GiftName:  userGift.Gift.Name,
		Month:     userGift.Month,
		Year:      userGift.Year,
		Status:    userGift.Status.String(),
	}
}

func FormatUserGifts(userGifts []*models.UserGift) []*UserGiftResponse {
	formattedUserGifts := []*UserGiftResponse{}
	for _, userGift := range userGifts {
		formattedUserGift := FormatUserGift(userGift)
		formattedUserGifts = append(formattedUserGifts, formattedUserGift)
	}
	return formattedUserGifts
}

package dtos

import (
	"time"

	"final-project-backend/internal/models"
)

type InvoiceResponse struct {
	ID               int64      `json:"id"`
	Code             string     `json:"code"`
	UserID           int64      `json:"user_id"`
	UserEmail        string     `json:"user_email"`
	Status           string     `json:"status"`
	Total            int        `json:"total"`
	OriginalPrice    int        `json:"original_price"`
	SubscriptionID   int64      `json:"subscription_id"`
	SubscriptionPlan string     `json:"subscription_plan"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	PurchasedAt      time.Time  `json:"purchased_at"`
}

type GetWaitingInvoiceByCodeProtectedResponse struct {
	ID             int64                `json:"id"`
	Code           string               `json:"code"`
	UserID         int64                `json:"user_id"`
	Status         string               `json:"status"`
	Total          int                  `json:"total"`
	SubscriptionID int64                `json:"subscription_id"`
	Subscription   SubscriptionResponse `json:"subscription"`
}

type GetAllInvoicesResponse = []*InvoiceResponse

type CreateInvoiceRequest struct {
	OriginalPrice  int    `json:"original_price" binding:"required"`
	SubscriptionID int    `json:"subscription_id" binding:"required"`
	VoucherCode    string `json:"voucher_code"`
}

type CreateInvoiceResponse = InvoiceResponse

type UpdateProcessedInvoiceRequest struct {
	IsSuccess *bool `json:"is_success" binding:"required"`
}

type UpdateProcessedInvoiceResponse struct {
	InvoiceResponse
	GiftsSent    []*GiftResponse    `json:"gifts_sent"`
	VouchersSent []*VoucherResponse `json:"vouchers_sent"`
}

func FormatInvoice(s *models.Invoice) *InvoiceResponse {
	return &InvoiceResponse{
		ID:               int64(s.ID),
		UserID:           s.UserID,
		Status:           s.Status.String(),
		Total:            s.Total,
		SubscriptionID:   s.SubscriptionID,
		PaidAt:           s.PaidAt,
		Code:             s.Code,
		UserEmail:        s.User.Email,
		SubscriptionPlan: s.Subscription.Name,
		PurchasedAt:      s.Model.CreatedAt,
		OriginalPrice:    s.OriginalPrice,
	}
}

func FormatInvoices(invoices []*models.Invoice) []*InvoiceResponse {
	formattedInvoices := []*InvoiceResponse{}
	for _, invoice := range invoices {
		formattedInvoice := FormatInvoice(invoice)
		formattedInvoices = append(formattedInvoices, formattedInvoice)
	}
	return formattedInvoices
}

package dtos

import "final-project-backend/internal/models"

type VoucherResponse struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Discount        int    `json:"discount"`
	MinimumSpending int    `json:"minimum_spending"`
}

type GetAllVouchersResponse = []*VoucherResponse

func FormatVoucher(voucher *models.Voucher) *VoucherResponse {
	return &VoucherResponse{
		ID:              int64(voucher.ID),
		Name:            voucher.Name,
		Discount:        voucher.Discount,
		MinimumSpending: voucher.MinimumSpending,
	}
}

func FormatVouchers(vouchers []*models.Voucher) []*VoucherResponse {
	formattedVouchers := []*VoucherResponse{}
	for _, voucher := range vouchers {
		formatted := FormatVoucher(voucher)
		formattedVouchers = append(formattedVouchers, formatted)
	}
	return formattedVouchers
}

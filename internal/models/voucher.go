package models

import "gorm.io/gorm"

type Voucher struct {
	gorm.Model
	Name            string `json:"name"`
	Discount        int    `json:"discount"`
	MinimumSpending int    `json:"minimum_spending"`
}

package models

import "gorm.io/gorm"

type Gift struct {
	gorm.Model
	Name            string `json:"name"`
	Stock           int    `json:"stock"`
	MinimumSpending int    `json:"minimum_spending"`
}

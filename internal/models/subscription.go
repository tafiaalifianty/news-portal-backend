package models

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	ID    int64  `gorm:"primary_key"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Quota int    `json:"quota"`
}

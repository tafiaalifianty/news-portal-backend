package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             int64          `gorm:"primary_key"`
	Email          string         `json:"email" gorm:"unique"`
	Password       string         `json:"password"`
	Role           Role           `json:"role"`
	Fullname       string         `json:"fullname"`
	Address        string         `json:"address"`
	Histories      []Post         `json:"history" gorm:"many2many:histories;"`
	Subscriptions  []Subscription `json:"subscriptions" gorm:"many2many:UserSubscriptions;"`
	ReferralCode   string         `json:"referral_code" gorm:"unique"`
	ReferredUserID int64          `json:"referred_user_id" gorm:"default:null"`
	UsersReferred  []User         `json:"users_referred" gorm:"foreignkey:referred_user_id"`
}

type Role string

const (
	Member Role = "member"
	Admin  Role = "admin"
)

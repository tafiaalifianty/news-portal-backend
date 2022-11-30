package models

type UserToken struct {
	UserID       int64 `gorm:"primaryKey"`
	Email        string
	Role         Role
	RefreshToken string
}

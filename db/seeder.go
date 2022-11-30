package db

import (
	"errors"

	"final-project-backend/internal/models"

	"gorm.io/gorm"
)

func seedAdmin(db *gorm.DB) {
	if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		db.Create(&models.User{
			Email:    "admin@email.com",
			Password: "$2a$10$yJzEUyzws/jaVeZKT4KKvOJe54twQZoL8USbb9v2fatCQkoxVsxfS",
			Role:     models.Admin,
			Fullname: "admin",
			Address:  "Address",
		})
	}
}

func seedCategories(db *gorm.DB) {
	if err := db.First(&models.Category{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var categories = []models.Category{{Name: "Food"}, {Name: "Technology"}, {Name: "Automotive"}, {Name: "Health"}, {Name: "Travel"}, {Name: "Finance"}, {Name: "Sport"}}
		db.Create(&categories)
	}
}

func seedPostTypes(db *gorm.DB) {
	if err := db.First(&models.PostType{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var types = []models.PostType{{Name: "VIP", Quota: 2}, {Name: "Premium", Quota: 1}, {Name: "Free", Quota: 0}}
		db.Create(&types)
	}
}

func seedGifts(db *gorm.DB) {
	if err := db.First(&models.Gift{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var categories = []models.Gift{{Name: "Coffee Brewing Tools", Stock: 10}, {Name: "Coffee", Stock: 10}, {Name: "Glass Mug", Stock: 10}}
		db.Create(&categories)
	}
}

func seedSubscriptions(db *gorm.DB) {
	if err := db.First(&models.Subscription{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		var categories = []models.Subscription{{Name: "Standard", Price: 30000, Quota: 5}, {Name: "Premium", Price: 50000, Quota: 10}, {Name: "Gold", Price: 90000, Quota: 20}}
		db.Create(&categories)
	}
}

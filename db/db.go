package db

import (
	"final-project-backend/config"
	"final-project-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	dsnConf := config.InitConfigDsn()

	var err error
	db, err = gorm.Open(postgres.Open(dsnConf), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.History{}, &models.UserSubscriptions{}, &models.Invoice{}, &models.UserToken{}, &models.Gift{}, &models.UserGift{}, &models.Voucher{}, &models.UserVoucher{}, &models.UserSpending{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&models.User{}, "Histories", &models.History{})
	if err != nil {
		return err
	}

	err = db.SetupJoinTable(&models.User{}, "Subscriptions", &models.UserSubscriptions{})
	if err != nil {
		return err
	}

	seedAdmin(db)
	seedCategories(db)
	seedGifts(db)
	seedPostTypes(db)
	seedSubscriptions(db)

	return nil

}

func Get() *gorm.DB {
	return db
}

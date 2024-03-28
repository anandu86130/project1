package database

import (
	"log"
	"os"
	"project1/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBconnect() {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	DB = db

	err = DB.AutoMigrate(&model.UserModel{}, &model.OTP{}, &model.AdminModel{}, &model.Category{}, &model.Product{}, &model.Address{}, &model.Cart{}, &model.Order{}, &model.Coupon{}, &model.Orderitems{}, &model.Paymentdetails{}, &model.Wallet{}, &model.Whishlist{})
	if err != nil {
		log.Fatal("failed to auto migrate", err)
	}
}

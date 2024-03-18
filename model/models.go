package model

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	UserID   uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password" gorm:"not null"`
	Status   bool   `gorm:"not null" json:"status"`
}

type OTP struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Otp   string `json:"otp"`
	Email string `json:"email" gorm:"not null;unique"`
	Exp   time.Time
}

type Address struct {
	AddressId uint `json:"address_id" gorm:"primaryKey;autoIncrement"`
	Address  string `json:"address" gorm:"not null;unique"`
	City     string `json:"city" gorm:"not null"`
	Landmark string `json:"landmark" gorm:"not null"`
	State    string `json:"state" gorm:"not null"`
	Country  string `json:"country" gorm:"not null"`
	Pincode  string `json:"pincode" gorm:"not null"`
	UserId   uint   `json:"userid" gorm:"not null"`
	User     UserModel
}

type Category struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null;unique"`
	Description string     `gorm:"not null" json:"description"`
	DeletedAt   *time.Time `gorm:"index"`
}

type Product struct {
	gorm.Model
	Product_name string     `json:"name"`
	ImagePath1   string     `json:"imagepath1"`
	ImagePath2   string     `json:"imagepath2"`
	ImagePath3   string     `json:"imagepath3"`
	Description  string     `json:"description"`
	Price        uint       `json:"price"`
	Size         string     `json:"size"`
	Quantity     uint       `json:"quantity"`
	DeletedAt    *time.Time `gorm:"index"`
}

type Cart struct {
	gorm.Model
	UserID    uint `gorm:"user_id"`
	User      UserModel
	ProductId uint `gorm:"product_id"`
	Product   Product
	Quantity  uint
}

type AdminModel struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

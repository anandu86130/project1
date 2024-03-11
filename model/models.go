package model

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"_id"`
	Name      string `json:"name" gorm:"not null"`
	Email     string `gorm:"unique;not null" json:"email"`
	Addresses string `json:"address" gorm:"not null"`
	Password  string `json:"password" gorm:"not null"`
	Status    bool   `gorm:"not null" json:"status"`
}

type OTP struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"_id"`
	Otp   string `json:"otp"`
	Email string `json:"email" gorm:"not null;unique"`
	Exp   time.Time
}
type Address struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"_id"`
	Address  string `json:"address" gorm:"not null"`
	City     string `json:"city" gorm:"not null"`
	Landmark string `json:"landmark" gorm:"not null"`
	State    string `json:"state" gorm:"not null"`
	Country  string `json:"country" gorm:"not null"`
	Pincode  string `json:"pincode" gorm:"not null"`
	User_ID  uint   `json:"userid" gorm:"not null"`
}

type Category struct {
	gorm.Model
	CategoryId  uint       `gorm:"primaryKey;autoIncrement" json:"_id"`
	Name        string     `json:"name" gorm:"not null;unique"`
	Description string     `gorm:"not null" json:"description"`
	DeletedAt   *time.Time `gorm:"index"`
}

type Product struct {
	gorm.Model
	ProductId    uint       `gorm:"primaryKey;autoIncrement" json:"_id"`
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

type AdminModel struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

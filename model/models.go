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
	Addres   string `json:"address"`
	City     string `json:"city"`
	Street   string `json:"street"`
	Landmark string `json:"landmark"`
	State    string `json:"state"`
	Country  string `json:"country"`
	Pincode  uint   `json:"pincode"`
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
	ProductId    uint   `gorm:"primaryKey;autoIncrement" json:"_id"`
	Product_name string `json:"name"`
	ImagePath1   string
	ImagePath2   string
	ImagePath3   string
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

package model

import "gorm.io/gorm"

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
	Email string `json:"email" gorm:"not null"`
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
	CategoryId uint   `gorm:"primaryKey;autoIncrement" json:"_id"`
	Name       string `json:"name" gorm:"not null;unique"`
	Products   uint   `gorm:"foreignKey:ProductId;not null" json:"products" `
}

type Product struct {
	ProductId    uint   `gorm:"primaryKey;autoIncrement" json:"_id"`
	Product_name string `json:"name"`
	Image        string `json:"image"`
	Price        uint   `json:"price"`
	Size         string `json:"size"`
	Quantity     uint   `json:"quantity"`
}

type AdminModel struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

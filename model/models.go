package model

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	UserID   uint   `gorm:"primaryKey" json:"user_id"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password" gorm:"not null"`
	Status   bool   `gorm:"not null" json:"status"`
}

type OTP struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Otp   string `json:"otp"`
	Email string `json:"email" gorm:"not null;unique"`
	Exp   time.Time
}

type Address struct {
	AddressId uint   `json:"address_id" gorm:"primaryKey"`
	Address   string `json:"address" gorm:"not null;unique"`
	City      string `json:"city" gorm:"not null"`
	Landmark  string `json:"landmark" gorm:"not null"`
	State     string `json:"state" gorm:"not null"`
	Country   string `json:"country" gorm:"not null"`
	Pincode   string `json:"pincode" gorm:"not null"`
	UserId    uint   `json:"userid" gorm:"not null"`
	User      UserModel
}

type Category struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null;unique"`
	Description string     `gorm:"not null" json:"description"`
	DeletedAt   *time.Time `gorm:"index"`
}

type Product struct {
	gorm.Model
	Product_name string `json:"name" gorm:"unique"`
	CategoryID   uint   `json:"category_id"`
	Category     Category
	ImagePath1   string     `json:"imagepath1"`
	ImagePath2   string     `json:"imagepath2"`
	ImagePath3   string     `json:"imagepath3"`
	Description  string     `json:"description"`
	Price        uint       `json:"price"`
	Size         string     `json:"size"`
	Quantity     uint       `json:"quantity"`
	DeletedAt    *time.Time `gorm:"index"`
}

type Coupon struct {
	gorm.Model
	Code      string    `gorm:"unique" json:"code"`
	Discount  float64   `json:"discount"`
	ValidFrom time.Time `json:"validfrom"`
	ValidTo   time.Time `json:"validto"`
}

type Cart struct {
	gorm.Model
	User      UserModel
	UserID    uint `json:"user_id"`
	Product   Product
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type Order struct {
	gorm.Model
	User          UserModel
	UserID        uint `json:"user_id"`
	CouponId      uint
	Code          string
	Totalquantity uint
	Totalamount   uint
	Paymentmethod string
	Address       Address
	AddressId     uint `json:"address_id"`
	Orderdate     time.Time
}

type Orderitems struct {
	gorm.Model
	Order             Order
	OrderID           uint `json:"order_id"`
	Product           Product
	ProductID         uint
	Quantity          uint
	Subtotal          uint
	Orderstatus       string `json:"orderstatus"`
	Ordercancelreason string
}

type Paymentdetails struct {
	gorm.Model
	PaymentId     string
	OrderId       string
	Reciept       uint
	Paymentstatus string
	Paymentamount int
}

type Wallet struct {
	gorm.Model
	UserID uint
	User   UserModel
	Amount float64
}

type Whishlist struct {
	gorm.Model
	UserID    uint
	User      UserModel
	ProductID uint
	Product   Product
}

type Rating struct {
	gorm.Model
	UserID    uint
	User      UserModel
	ProductID uint
	Product   Product
	Rating    uint
	Review    string
}
type Productoffer struct {
	gorm.Model
	ProductID uint
	Product   Product
	Offer     uint
}
type AdminModel struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

package routes

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	userid := c.GetUint("userid")
	var cart []model.Cart
	result := database.DB.Preload("Product").Where("user_id=?", userid).First(&cart)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}
	var totalamount uint
	var totalquantity uint
	var productid uint
	var price uint
	for _, cartitem := range cart {
		price := cartitem.Quantity * cartitem.Product.Price
		totalamount += price
		totalquantity += cartitem.Quantity
		productid = cartitem.ProductID
		price = cartitem.Product.Price
	}
	addressIDStr := c.PostForm("address_id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid addressid"})
		return
	}
	order := model.Order{
		UserID:        userid,
		ProductID:     productid,
		Totalamount:   totalamount,
		Price:         price,
		Totalquantity: totalquantity,
		Paymentmethod: "COD",
		AddressId:     uint(addressID),
	}

	if result := database.DB.Create(&order).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}
}

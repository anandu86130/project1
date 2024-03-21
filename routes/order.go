package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func Adminorderview(c *gin.Context) {
	var order []model.Order
	if result := database.DB.Preload("Order").Find(&order).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find order"})
		return
	}

	for _, orderlist := range order {
		c.JSON(http.StatusOK, gin.H{
			"order id":       orderlist.ID,
			"total amount":   orderlist.Totalamount,
			"user id":        orderlist.UserID,
			"payment method": orderlist.Paymentmethod,
			"order date":     orderlist.Orderdate,
			"address id":     orderlist.AddressId,
		})
	}
}


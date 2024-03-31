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
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
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

func Adminorderstatus(c *gin.Context) {
	orderid := c.Param("ID")
	var orderitem model.Orderitems
	if err := database.DB.Find(&orderitem, orderid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
	}
	status := c.Request.FormValue("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Please provide status"})
		return
	}
	if err := database.DB.Model(&orderitem).Update("orderstatus", status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Status changed successfully"})

}

func Admincancelorders(c *gin.Context) {
	orderid := c.Param("ID")
	if orderid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Please provide the order id"})
		return
	}
	var order model.Orderitems
	if result := database.DB.Find(&order, orderid).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Cant find the order"})
		return
	}

	if order.Orderstatus == "cancelled" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Order already cancelled"})
		return
	} else {
		order.Orderstatus = "cancelled"
		if result := database.DB.Save(&order.Orderstatus).Error; result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update order status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Message": "Order cancelled successfully"})
	}
}

package routes

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	userid := c.GetUint("userid")
	var cart []model.Cart
	result := database.DB.Preload("Product").Where("user_id=?", userid).Find(&cart)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	addressIDStr := c.Param("ID")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid addressid"})
		return
	}

	var coupon model.Coupon
	couponcode := c.Request.FormValue("code")
	if couponcode != "" {
		if result := database.DB.Where("code=?", couponcode).First(&coupon).Error; result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid coupon code"})
			return
		}

		currenttime := time.Now()
		if currenttime.Before(coupon.ValidFrom) || currenttime.After(coupon.ValidTo) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "coupon is not valid"})
			return
		}
	}
	// var orders []model.Order
	var orderitems []model.Orderitems
	var totalamount uint
	for _, cartitem := range cart {
		totalamount += cartitem.Quantity * cartitem.Product.Price

		orderlist := model.Orderitems{
			OrderID:           cartitem.ID,
			ProductID:         cartitem.ProductID,
			Quantity:          cartitem.Quantity,
			Subtotal:          cartitem.Quantity * cartitem.Product.Price,
			Orderstatus:       "pending",
			Ordercancelreason: "",
		}
		orderitems = append(orderitems, orderlist)

		if err := database.DB.Model(&cartitem.Product).Update("quantity", cartitem.Product.Quantity-cartitem.Quantity).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "failed to update quantity"})
			return
		}
	}
	currenttime := time.Now()

	totalamount = uint(float64(totalamount) - coupon.Discount)

	order := model.Order{
		UserID:        userid,
		CouponId:      coupon.ID,
		Code:          coupon.Code,
		Totalquantity: uint(len(cart)),
		Totalamount:   totalamount,
		Paymentmethod: "COD",
		AddressId:     uint(addressID),
		Orderdate:     currenttime,
	}

	if result := database.DB.Create(&order).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	for i := range orderitems {
		orderitems[i].OrderID = order.ID
	}
	if create := database.DB.Create(&orderitems).Error; create != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create orderlist"})
		return
	}

	if err := database.DB.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "order placed successfully",
		"Totalamount": totalamount,
	})
}

func Orderview(c *gin.Context) {
	var order []model.Order
	userid := c.GetUint("userid")
	if result := database.DB.Where("user_id=?", userid).Find(&order).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	for _, orderdetails := range order {
		c.JSON(http.StatusOK, gin.H{
			"order id":       orderdetails.ID,
			"Amount":         orderdetails.Totalamount,
			"payment method": orderdetails.Paymentmethod,
			"order date":     orderdetails.Orderdate,
		})
	}
}

func Orderdetails(c *gin.Context) {
	var orderlist []model.Orderitems
	orderid := c.Param("ID")
	if result := database.DB.Where("order_id", orderid).Preload("Order").Preload("Product").Find(&orderlist).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find order details"})
		return
	}
	for _, orderitem := range orderlist {
		c.JSON(http.StatusOK, gin.H{
			"order item id": orderitem.ID,
			"Product":       orderitem.ProductID,
			"product name":  orderitem.Product.Product_name,
			"order date":    orderitem.Order.Orderdate,
			"Amount":        orderitem.Subtotal,
			"Quantity":      orderitem.Quantity,
			"Status":        orderitem.Orderstatus,
			"AddressID":     orderitem.Order.AddressId,
		})
	}
}

func Cancelorder(c *gin.Context) {
	var orderlist model.Orderitems
	var productQuantity model.Product
	var order model.Order
	orderitemid := c.Param("ID")
	if orderitemid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "please give the orderid"})
		return
	}

	orderID, err := strconv.ParseUint(orderitemid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid order ID"})
		return
	}

	if result := database.DB.First(&orderlist, orderID).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find order"})
		return
	}
	reason := c.Request.FormValue("reason")
	if reason == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "please give the reason"})
		return
	}

	if orderlist.Orderstatus == "cancelled" {
		c.JSON(http.StatusOK, gin.H{"message": "order already cancelled"})
		return
	}

	orderlist.Orderstatus = "cancelled"
	orderlist.Ordercancelreason = reason

	if err := database.DB.Save(&orderlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
		return
	}

	database.DB.First(&productQuantity, orderlist.ProductID)
	productQuantity.Quantity += orderlist.Quantity
	if err := database.DB.Save(&productQuantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order successfully cancelled"})

	cancelledamount := orderlist.Subtotal
	database.DB.Model(&order).Updates(model.Order{
		Totalamount: order.Totalamount - uint(cancelledamount),
	})

}

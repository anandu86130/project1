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
	result := database.DB.Preload("Product").Where("user_id=?", userid).First(&cart)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}
	var totalamount uint
	var totalquantity uint
	// var productid uint
	var Productprice uint
	for _, cartitem := range cart {
		price := cartitem.Quantity * cartitem.Product.Price
		totalamount += price
		totalquantity += cartitem.Quantity
		// productid = cartitem.ProductID
		Productprice = cartitem.Product.Price

		if err := database.DB.Model(&cartitem.Product).Update("quantity", cartitem.Product.Quantity-cartitem.Quantity).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "failed to update quantity"})
			return
		}
	}

	couponcode := c.Request.FormValue("code")
	var coupon model.Coupon
	if result := database.DB.Where("code=?", couponcode).First(&coupon).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid coupon code"})
		return
	}

	currenttime := time.Now()
	if currenttime.Before(coupon.ValidFrom) || currenttime.After(coupon.ValidTo) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "coupon is not valid"})
		return
	}

	totalamount = uint(float64(totalamount) - coupon.Discount)

	addressIDStr := c.Param("address_id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid addressid"})
		return
	}

	var orders []model.Order
	for _, val := range cart {
		order := model.Order{
			UserID:        userid,
			ProductID:     val.ProductID,
			Totalamount:   totalamount,
			CouponId:      coupon.ID,
			Code:          coupon.Code,
			Price:         Productprice,
			Totalquantity: totalquantity,
			Orderdate:     currenttime,
			Paymentmethod: "COD",
			AddressId:     uint(addressID),
		}
		orders = append(orders, order)

		orderlist := model.Orderitems{
			OrderID:           order.ID,
			ProductID:         order.ProductID,
			Quantity:          order.Totalquantity,
			Subtotal:          order.Totalamount,
			Orderstatus:       "pending",
			Ordercancelreason: "",
		}

		if create := database.DB.Create(&orderlist).Error; create != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create orderlist"})
			return
		}
	}

	for _, orderitem := range orders {
		if result := database.DB.Create(&orderitem).Error; result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
			return
		}
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
	if result := database.DB.Preload("Orderitems").Find(&order).Where("user_id=?", userid).Error; result != nil {
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

// func Cancelorder(c *gin.Context) {
// 	var orderlist model.Orderitems
// 	var productQuantity model.Product
// 	orderitemid := c.Param("ID")
// 	if orderitemid == "" {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "please give the orderid"})
// 		return
// 	}
// 	reason := c.Request.FormValue("reason")
// 	if reason == "" {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "please give the reason"})
// 		return
// 	}

// 	if result := database.DB.Where("order_id=?", orderitemid).First(&orderlist).Error; result != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find order"})
// 		return
// 	}

// 	if orderlist.Orderstatus == "cancelled" {
// 		c.JSON(http.StatusOK, gin.H{"message": "order already cancelled"})
// 		return
// 	}

// 	orderlist.Orderstatus = "cancelled"
// 	orderlist.Ordercancelreason = reason

// 	if err := database.DB.Save(&orderlist).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
// 		return
// 	}

// 	database.DB.First(&productQuantity, orderlist.ProductID)
// 	productQuantity.Quantity += orderlist.Quantity
// 	if err := database.DB.Save(&productQuantity).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add quantity"})
// 		return
// 	}
// 	var orderamount model.Order
// 	if err := database.DB.First(&orderamount, orderlist.OrderID).Error; err != nil{
// 		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to find order details"})
// 		return
// 	}
// 	var couponremove model.Coupon
// 	if orderamount.
// }

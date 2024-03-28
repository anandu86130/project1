package routes

import (
	"math/rand"
	"net/http"
	"project1/database"
	"project1/model"
	"project1/payment"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	userid := c.GetUint("userid")
	var cart []model.Cart
	result := database.DB.Preload("Product").Where("user_id=?", userid).Find(&cart)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find user"})
		return
	}

	addressIDStr := c.Param("ID")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid addressid"})
		return
	}

	var coupon model.Coupon
	couponcode := c.Request.FormValue("code")
	if couponcode != "" {
		if result := database.DB.Where("code=?", couponcode).First(&coupon).Error; result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Invalid coupon code"})
			return
		}

		currenttime := time.Now()
		if currenttime.Before(coupon.ValidFrom) || currenttime.After(coupon.ValidTo) {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Coupon is not valid"})
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
			c.JSON(http.StatusOK, gin.H{"Error": "Failed to update quantity"})
			return
		}
	}
	currenttime := time.Now()

	totalamount = uint(float64(totalamount) - coupon.Discount)

	paymentmethod := c.Request.FormValue("paymentmethod")
	if paymentmethod == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give paymentmethod PAYNOW or COD"})
		return
	}
	if paymentmethod == "PAYNOW" {

		razorId, errr := payment.Paymenthandler(strconv.Itoa(int(orderitems[0].ID)), int(totalamount))
		if errr != nil {
			c.JSON(402, gin.H{"Error": errr})
			return
		}
		recieptID := generateReceiptID()
		create := model.Paymentdetails{
			OrderId:       razorId,
			Paymentamount: int(totalamount),
			Reciept:       uint(recieptID),
			Paymentstatus: "Failed",
		}
		database.DB.Create(&create)
		c.JSON(200, gin.H{"Message": "Complete the payment", "Order": razorId})
	}

	order := model.Order{
		UserID:        userid,
		CouponId:      coupon.ID,
		Code:          coupon.Code,
		Totalquantity: uint(len(cart)),
		Totalamount:   totalamount,
		Paymentmethod: paymentmethod,
		AddressId:     uint(addressID),
		Orderdate:     currenttime,
	}

	if result := database.DB.Create(&order).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create order"})
		return
	}

	for i := range orderitems {
		orderitems[i].OrderID = order.ID
	}
	if create := database.DB.Create(&orderitems).Error; create != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create orderlist"})
		return
	}

	if err := database.DB.Delete(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message":     "Order placed successfully",
		"Totalamount": totalamount,
	})
}

func generateReceiptID() int64 {
	timestamp := time.Now().Unix()
	return timestamp*1000 + rand.Int63n(1000)
}

func Orderview(c *gin.Context) {
	var order []model.Order
	userid := c.GetUint("userid")
	if result := database.DB.Where("user_id=?", userid).Find(&order).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"Eror": "Failed to find order details"})
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
	// var order model.Order
	orderitemid := c.Param("ID")
	if orderitemid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give the orderid"})
		return
	}

	orderID, err := strconv.ParseUint(orderitemid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Invalid order ID"})
		return
	}

	if result := database.DB.First(&orderlist, orderID).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
		return
	}
	reason := c.Request.FormValue("reason")
	if reason == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give the reason"})
		return
	}

	if orderlist.Orderstatus == "cancelled" {
		c.JSON(http.StatusOK, gin.H{"Message": "Order already cancelled"})
		return
	}

	if err := database.DB.Save(&orderlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update order"})
		return
	}

	var orderamount model.Order
	if err := database.DB.First(&orderamount, orderlist.OrderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order details"})
		return
	}

	database.DB.First(&productQuantity, orderlist.ProductID)
	productQuantity.Quantity += orderlist.Quantity
	if err := database.DB.Save(&productQuantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add quantity"})
		return
	}

	var coupon model.Coupon
	if orderamount.Code != "" {
		if err := database.DB.First(&coupon, "code=?", orderamount.Code).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Cant find coupon code"})
			return
		} else {
			orderamount.Totalamount += uint(coupon.Discount)
			orderamount.Totalamount -= orderlist.Subtotal
			orderamount.Code = ""
		}
		if err := database.DB.Save(&orderamount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update order details"})
			return
		}
	}

	if orderlist.Order.Paymentmethod != "PAYNOW" {
		orderlist.Orderstatus = "cancelled"
		orderlist.Ordercancelreason = reason
		if err := database.DB.Save(&orderlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update order item"})
		}
	}

	var wallet model.Wallet
	userID := c.GetUint("userid")
	if err := database.DB.First(&wallet, "user_id=?", userID).Error; err != nil {
		wallet = model.Wallet{
			UserID: userID,
			Amount: 0,
		}
		if err := database.DB.Create(&wallet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create wallet"})
			return
		}
	}
	wallet.Amount += float64(orderlist.Subtotal)

	if err := database.DB.Save(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save amount to wallet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Order successfully cancelled"})
}

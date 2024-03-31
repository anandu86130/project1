package routes

import (
	"fmt"
	"net/http"
	"project1/database"
	"project1/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CartView(c *gin.Context) {
	var cart []model.Cart
	var show []gin.H
	userID := c.GetUint("userid")
	var totalamount = 0
	var count = 0
	err := database.DB.Preload("Product").Where("user_id=?", userID).Find(&cart)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find cart"})
		return
	}

	for _, val := range cart {
		show = append(show, gin.H{
			"product name":     val.Product.Product_name,
			"product image":    val.Product.ImagePath1,
			"product quantity": val.Quantity,
			"product price":    val.Product.Price,
			"product id":       val.Product.ID,
		})
		price := int(val.Quantity) * int(val.Product.Price)
		totalamount += price
		count += 1
	}
	if totalamount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "No products added to cart"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"cart":           show,
			"total products": count,
			"total Amount":   totalamount,
		})
	}
	cart = []model.Cart{}
}

func Addtocart(c *gin.Context) {
	userID := c.GetUint("userid")
	idStr := c.Param("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to convert"})
		return
	}

	var product model.Product
	if err := database.DB.Where("id=? AND deleted_at is NULL", id).First(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}

	var cart model.Cart
	result := database.DB.Where("user_id = ? AND Product_id = ?", userID, id).First(&cart)
	if result.Error != nil {
		cart = model.Cart{
			UserID:    userID,
			ProductID: uint(id),
			Quantity:  1,
		}
		fmt.Println("userid===========================================================================================", userID)
		fmt.Println("productid===========================================================================================", id)
		if err := database.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save to cart"})
			return
		}
	} else {
		cart.Quantity++
		database.DB.Save(&cart)
	}
	fmt.Println("userid===========================================================================================", userID)
	fmt.Println("productid===========================================================================================", id)
	c.JSON(http.StatusOK, gin.H{"Message": "Product added successfully"})
}

func Deletecart(c *gin.Context) {
	userID := c.GetUint("userid")
	idStr := c.Param("ID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to convert"})
		return
	}
	var cart model.Cart
	result := database.DB.Where("user_id=? AND product_id=?", userID, id).First(&cart)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}

	delete := database.DB.Delete(&cart)
	if delete.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete product from cart"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "Successfully deleted product from the cart"})
	}
}

package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func Coupon(c *gin.Context) {
}

func Addcoupon(c *gin.Context) {
	var coupon model.Coupon
	err := c.ShouldBindJSON(&coupon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind"})
		return
	}

	if result := database.DB.Create(&coupon).Error; result != nil {
		c.JSON(http.StatusOK, gin.H{"error": "failed to create coupon"})
		return
	}else{
		c.JSON(http.StatusOK, gin.H{"message": "coupon created successfully"})
}
}

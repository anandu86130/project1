package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func Coupon(c *gin.Context) {
	var coupon []model.Coupon
	if result := database.DB.Find(&coupon).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find coupon"})
		return
	}

	for _, val := range coupon {
		c.JSON(http.StatusOK, gin.H{
			"coupon Id":         val.ID,
			"coupon code":       val.Code,
			"coupon discount":   val.Discount,
			"coupon valid from": val.ValidFrom,
			"coupon valid to":   val.ValidTo,
		})
	}
	coupon = []model.Coupon{}
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
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "coupon created successfully"})
	}
}

func Deletecoupon(c *gin.Context) {
	var deletecoupon model.Coupon
	id := c.Param("ID")
	if result := database.DB.Where("id=?", id).Delete(&deletecoupon).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete coupon"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "coupon deleted successfully"})
}

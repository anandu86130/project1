package routes

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Whishlist(c *gin.Context) {
	userid := c.GetUint("userid")
	var whishlist []model.Whishlist
	if err := database.DB.Preload("Product").Where("user_id=?", userid).Find(&whishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find whishlist"})
		return
	}

	var products []gin.H
	for _, val := range whishlist {
		details := gin.H{
			"whishlistid":   val.ID,
			"product name":  val.Product.Product_name,
			"product price": val.Product.Price,
		}
		products = append(products, details)
	}
	c.JSON(http.StatusOK, products)
}

func Addtowhishlist(c *gin.Context) {
	userid := c.GetUint("userid")
	id := c.Param("ID")
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give product id"})
		return
	}
	productid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Invalid product id"})
		return
	}
	var product model.Product
	if err := database.DB.First(&product, productid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}
	var whishlist model.Whishlist
	if err := database.DB.Where("user_id=? AND product_id=?", userid, productid).First(&whishlist).Error; err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Product already added to whishlist"})
		return
	}

	whishlist = model.Whishlist{
		UserID:    userid,
		ProductID: uint(productid),
	}
	if err := database.DB.Create(&whishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add to whishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Product added to whishlist"})
}

func Deletewhishlist(c *gin.Context) {
	id := c.Param("ID")
	whishlistid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find whishlist"})
		return
	}
	var whishlist model.Whishlist
	if err := database.DB.Where("id=?", whishlistid).First(&whishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find whishlist"})
		return
	}
	if err := database.DB.Delete(&whishlist); err != nil {
		c.JSON(http.StatusOK, gin.H{"Message": "Whishlist deleted successfully"})
	}
}

package routes

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Addrating(c *gin.Context) {
	var rating model.Rating
	err := c.BindJSON(&rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find rating"})
		return
	}
	userid := c.GetUint("userid")

	id := c.Param("ID")
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give product id"})
		return
	}
	productid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch id"})
		return
	}
	var product model.Product
	if err := database.DB.First(&product, productid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product id"})
		return
	}
	var ratings model.Rating
	if err := database.DB.Where("user_id = ? AND product_id = ?", userid, productid).First(&ratings).Error; err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Rating for this product already exist"})
		return
	}
	if rating.Rating > 5 {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please provide rating out of 5"})
		return
	}
	if rating.Review == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please provide the review"})
		return
	}
	rating = model.Rating{
		UserID:    userid,
		ProductID: uint(productid),
		Rating:    rating.Rating,
		Review:    rating.Review,
	}
	if err := database.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create rating"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Rating added successfully"})
}

package routes

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func Categoryfilter(c *gin.Context) {
	filter := c.Request.FormValue("category")
	search := strings.ToLower(filter)
	var category model.Category
	if err := database.DB.Where("name=?", search).First(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find category"})
		return
	}
	var product []model.Product
	if err := database.DB.Preload("category").Where("category_id=?", category.ID).Find(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": product})
}

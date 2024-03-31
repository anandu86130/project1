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
	var details []gin.H
	var product []model.Product
	if err := database.DB.Where("category_id=?", category.ID).Find(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}
	for _, v := range product {
		details = append(details, gin.H{
			"Name":        v.Product_name,
			"Imagepath1":  v.ImagePath1,
			"Imagepath2":  v.ImagePath2,
			"Imagepath3":  v.ImagePath3,
			"Description": v.Description,
			"Price":       v.Price,
			"Size":        v.Size,
			"Quantity":    v.Quantity,
		})
	}
	c.JSON(http.StatusOK, gin.H{"Message": details})
}

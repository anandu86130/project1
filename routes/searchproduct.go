package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func Productsearch(c *gin.Context) {
	search := c.Request.FormValue("search")
	if search == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please send the search"})	
		return
	}
	var details []gin.H
	var product []model.Product
	if result := database.DB.Where("product_name ILIKE ?", "%"+search+"%").Find(&product).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}
	for _, v := range product {
		details = append(details, gin.H{
			"ProductID":   v.ID,
			"ProductName": v.Product_name,
			"CategoryID":  v.CategoryID,
			"ImagePath1":  v.ImagePath1,
			"ImagePath2":  v.ImagePath2,
			"ImagePath3":  v.ImagePath3,
			"Description": v.Description,
			"Price":       v.Price,
			"Size":        v.Size,
			"Quantity":    v.Quantity,
		})
	}
	c.JSON(http.StatusOK, details)
}

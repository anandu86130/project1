package bestcategory

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func BestSellingCategory(c *gin.Context) {
	var orderItems []model.Orderitems

	if err := database.DB.Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
		return
	}

	productQuantity := make(map[uint]uint)

	var bestSellingProduct model.Product
	var maxQuantity uint

	for _, item := range orderItems {
		productID := item.ProductID
		quantity := item.Quantity

		productQuantity[productID] += quantity

		for productID, quantity := range productQuantity {
			if quantity > maxQuantity {
				maxQuantity = quantity

				if err := database.DB.Preload("Category").First(&bestSellingProduct, productID).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find best selling product"})
					return
				}
			}
		}
	}
	var details []gin.H
	details = append(details, gin.H{
		"categoryName": bestSellingProduct.Category.Name,
		"Description":  bestSellingProduct.Category.Description,
	})
	c.JSON(http.StatusOK, gin.H{"Message": details})
}

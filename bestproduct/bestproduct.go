package bestproduct

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func BestSellingProduct(c *gin.Context) {
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

				if err := database.DB.First(&bestSellingProduct, productID).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find best selling product"})
					return
				}
			}
		}
	}
	var details []gin.H
	details = append(details, gin.H{
		"BestSellingProduct": bestSellingProduct.Product_name,
		"ProductPrice":       bestSellingProduct.Price,
		"ImagePath1":         bestSellingProduct.ImagePath1,
		"ImagePath2":         bestSellingProduct.ImagePath2,
		"ImagePath3":         bestSellingProduct.ImagePath3,
		"Description":        bestSellingProduct.Description,
		"Size":               bestSellingProduct.Size,
		"TotalQuantitySold":  maxQuantity,
	})
	c.JSON(http.StatusOK, gin.H{"Message": details})
}

package bestcategory

import (
	"net/http"
	"project1/database"
	"project1/model"
	"sort"

	"github.com/gin-gonic/gin"
)

func BestSellingCategory(c *gin.Context) {
	var orderItems []model.Orderitems

	if err := database.DB.Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
		return
	}

	categoryQuantity := make(map[uint]uint)

	for _, item := range orderItems {
		productID := item.ProductID

		var product model.Product
		if err := database.DB.Preload("Category").First(&product, productID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
			return
		}

		categoryID := product.CategoryID
		quantity := item.Quantity

		categoryQuantity[categoryID] += quantity
	}

	type CategoryQuantity struct {
		CategoryID uint
		Quantity   uint
	}

	var categoryQuantities []CategoryQuantity
	for categoryID, quantity := range categoryQuantity {
		categoryQuantities = append(categoryQuantities, CategoryQuantity{CategoryID: categoryID, Quantity: quantity})
	}

	sort.Slice(categoryQuantities, func(i, j int) bool {
		return categoryQuantities[i].Quantity > categoryQuantities[j].Quantity
	})

	var topCategories []model.Category
	for i := 0; i < 10 && i < len(categoryQuantities); i++ {
		var category model.Category
		if err := database.DB.First(&category, categoryQuantities[i].CategoryID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find category"})
			return
		}
		topCategories = append(topCategories, category)
	}
	var details []gin.H
	for _, category := range topCategories {
		details = append(details, gin.H{
			"categoryName":      category.Name,
			"Description":       category.Description,
			"TotalQuantitySold": categoryQuantity[category.ID],
		})
	}
	c.JSON(http.StatusOK, gin.H{"Message": details})
}

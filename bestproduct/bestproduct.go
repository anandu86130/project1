package bestproduct

import (
	"net/http"
	"project1/database"
	"project1/model"
	"sort"

	"github.com/gin-gonic/gin"
)

func BestSellingProduct(c *gin.Context) {
	var orderItems []model.Orderitems

	if err := database.DB.Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find order"})
		return
	}

	productQuantity := make(map[uint]uint)

	for _, item := range orderItems {
		productQuantity[item.ProductID] += item.Quantity
	}

	type ProductQuantity struct {
		ProductID uint
		Quantity  uint
	}

	var productQuantities []ProductQuantity
	for productID, quantity := range productQuantity {
		productQuantities = append(productQuantities, ProductQuantity{ProductID: productID, Quantity: quantity})
	}

	sort.Slice(productQuantities, func(i, j int) bool {
		return productQuantities[i].Quantity > productQuantities[j].Quantity
	})

	var topProducts []model.Product
	for i:=0; i<10 && i<len(productQuantities); i++{
		var product model.Product
		if err := database.DB.First(&product, productQuantities[i].ProductID).Error; err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"Error":"Failed to find product details"})
			return
		}
		topProducts = append(topProducts, product)
	}
	var details []gin.H
	for _, product := range topProducts{
		details = append(details, gin.H{
			"BestSellingProduct": product.Product_name,
			"ProductPrice":       product.Price,
			"ImagePath1":         product.ImagePath1,
			"ImagePath2":         product.ImagePath2,
			"ImagePath3":         product.ImagePath3,
			"Description":        product.Description,
			"Size":               product.Size,
			"TotalQuantitySold":  productQuantity[product.ID],
		})
	}
	c.JSON(http.StatusOK, gin.H{"Message": details})
}

package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

var Offeramount uint

func Addproductoffer(c *gin.Context) {
	var offer model.Productoffer
	err := c.BindJSON(&offer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind JSON"})
		return
	}
	var offers model.Productoffer
	if err := database.DB.Where("product_id=?", offer.ID).First(&offers).Error; err != nil {
		Offeramount = offer.Offer
		// product.Price -= offer.Price
		productoffer := model.Productoffer{
			ProductID: offer.ID,
			Offer:     Offeramount,
		}
		database.DB.Create(&productoffer)
		c.JSON(http.StatusOK, gin.H{"Message": "Prouct offer added successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Offer for this product already exists"})
		return
	}
	// var product model.Product
	// if err := database.DB.Where("id=?", offer.ID).First(&product).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product id"})
	// 	return
	// }

	// fmt.Println("=======================================================================================================", productoffer.ProductID)

	// if err := database.DB.Save(&product).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add product offer"})
	// 	return
	// }
}

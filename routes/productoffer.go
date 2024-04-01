package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

var Offeramount uint

func Addproductoffer(c *gin.Context) {
	var offer model.Product
	err := c.BindJSON(&offer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind JSON"})
		return
	}
	var product model.Product
	if err := database.DB.First(&product, offer.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product id"})
		return
	}
	Offeramount = offer.Price
	product.Price -= offer.Price
	productoffer := model.Productoffer{
		ProductID: product.ID,
		Offer:     offer.Price,
	}
	if err := database.DB.Model(&productoffer).Updates(productoffer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create offer"})
		return
	}
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to add product offer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Prouct offer added successfully"})
}

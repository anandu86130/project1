package routes

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func Sortproduct(c *gin.Context) {
	search := c.Request.FormValue("sort")
	sort := strings.ToLower(search)
	var product []model.Product
	switch sort {
	case "asc":
		database.DB.Order("product_name asc").Find(&product)
	case "desc":
		database.DB.Order("product_name desc").Find(&product)
	case "high to low":
		database.DB.Order("price desc").Find(&product)
	case "low to high":
		database.DB.Order("price asc").Find(&product)
	case "new arrivals":
		database.DB.Order("created_at desc").Find(&product)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please give correct sort"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": product})
}

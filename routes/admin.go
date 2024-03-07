package routes

import (
	"net/http"
	"project1/database"
	"project1/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func Signin(c *gin.Context) {
	var admin model.AdminModel
	err := c.ShouldBindJSON(&admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	var check model.AdminModel
	user := database.DB.Where("email=?", &admin.Email).Find(&check)
	if user.Error != nil {
		c.JSON(http.StatusUnauthorized, "invalid email or password")
		return
	}
	if admin.Password != check.Password {
		c.JSON(http.StatusUnauthorized, "invalid email or password")
		return
	}

	c.JSON(http.StatusOK, "successfully signed in")

}

func Getuser(c *gin.Context) {
	var users []model.UserModel
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, "cant load user")
		return
	}
	var responseData []gin.H
	for _, user := range users {
		userdata := gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
		responseData = append(responseData, userdata)
	}
	c.JSON(http.StatusOK, responseData)
}

// category
func Category(c *gin.Context) {
	var category []model.Category
	resutlt := database.DB.Find(&category)
	if resutlt.Error != nil {
		c.JSON(http.StatusBadRequest, "failed to load category")
		return
	}

	var categories []gin.H
	for _, categorydetails := range category {
		fetchedcategory := gin.H{
			"id":       categorydetails.CategoryId,
			"name":     categorydetails.Name,
			"products": categorydetails.Products,
		}
		categories = append(categories, fetchedcategory)
	}
	c.JSON(http.StatusOK, categories)
}

func Addcategory(c *gin.Context) {
	var category model.Category
	err := c.BindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind category")
		return
	}

	category.Name = strings.ToLower(category.Name)

	var dbcat model.Category
	result := database.DB.Where("LOWER(name) = ?", category.Name).First(&dbcat)
	if result.Error == nil {
		c.JSON(http.StatusConflict, "this category already exists")
		return
	}

	create := database.DB.Create(&category)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to create category")
		return
	}

	c.JSON(http.StatusOK, "category created successfully")
}

func Aproduct(c *gin.Context) {
	var product []model.Product
	result := database.DB.Find(&product)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, "failed to connect products")
		return
	}

	var products []gin.H
	for _, fetchedproducts := range product {
		productdetails := gin.H{
			"id":       fetchedproducts.ProductId,
			"name":     fetchedproducts.Product_name,
			"image":    fetchedproducts.Image,
			"price":    fetchedproducts.Price,
			"size":     fetchedproducts.Size,
			"quantity": fetchedproducts.Quantity,
		}
		products = append(products, productdetails)
	}
	c.JSON(http.StatusOK, products)
}

func Addproduct(c *gin.Context) {
	var product model.Product
	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	product.Product_name = strings.ToLower(product.Product_name)

	var dbproduct model.Product
	database.DB.Where("LOWER(product_name) = ?", product.Product_name).Find(&dbproduct)
	if dbproduct.Product_name == product.Product_name {
		c.JSON(http.StatusConflict, "this product already exists")
		return
	}

	create := database.DB.Create(&product)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to create product")
		return
	}

	c.JSON(http.StatusOK, "product created successfully")
}

func Blockuser(c *gin.Context) {
	var status model.UserModel
	id := c.Param("ID")
	database.DB.First(&status,id)
	if status.Status {
		database.DB.Model(&status).Update("status",false)
		c.JSON(200,"User Blocked")
	}else{
		database.DB.Model(&status).Update("status",true)
		c.JSON(200,"User Unblocked")
	}
}

func Editcategory(c *gin.Context) {

}

func Deletecategory(c *gin.Context) {

}

func Editproduct(c *gin.Context) {

}

func Deleteproduct(c *gin.Context) {

}

package routes

import (
	"net/http"
	"project1/database"
	"project1/jwt"
	"project1/model"
	"time"

	"github.com/gin-gonic/gin"
)

var Product model.Product

const RoleAdmin = "admin"

func Signin(c *gin.Context) {
	var admin model.AdminModel
	err := c.ShouldBindJSON(&admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind"})
		return
	}

	var check model.AdminModel
	user := database.DB.Where("email=?", &admin.Email).Find(&check)
	if user.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if admin.Password != check.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	jwt.JwtToken(c, admin.ID, admin.Email, RoleAdmin)
	c.JSON(http.StatusOK, gin.H{"message": "successfully signed in"})

}

func Getuser(c *gin.Context) {
	var users []model.UserModel
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant load user"})
		return
	}
	var responseData []gin.H
	for _, user := range users {
		userdata := gin.H{
			"Id":    user.UserID,
			"Name":  user.Name,
			"Email": user.Email,
		}
		responseData = append(responseData, userdata)
	}
	c.JSON(http.StatusOK, responseData)
}

// category
func Category(c *gin.Context) {
	var category []model.Category
	result := database.DB.Find(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load category"})
		return
	}

	var categories []gin.H
	for _, categorydetails := range category {
		fetchedcategory := gin.H{
			"Categoryid":  categorydetails.ID,
			"Name":        categorydetails.Name,
			"Description": categorydetails.Description,
		}
		categories = append(categories, fetchedcategory)
	}
	c.JSON(http.StatusOK, categories)
}

func Addcategory(c *gin.Context) {
	var category model.Category
	err := c.BindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind category"})
		return
	}
	var dbcat model.Category
	database.DB.Where("name = ?", category.Name).First(&dbcat)
	if dbcat.Name == category.Name {
		c.JSON(http.StatusConflict, gin.H{"error": "this category already exists"})
		return
	}

	create := database.DB.Create(&category)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category created successfully"})
}

func Aproduct(c *gin.Context) {
	var product []model.Product
	result := database.DB.Preload("Category").Find(&product)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to connect products"})
		return
	}
	var products []gin.H
	for _, fetchedproducts := range product {
		productdetails := gin.H{
			"Productid":   fetchedproducts.ID,
			"Name":        fetchedproducts.Product_name,
			"Imagepath1":  fetchedproducts.ImagePath1,
			"Imagepath2":  fetchedproducts.ImagePath2,
			"Imagepath3":  fetchedproducts.ImagePath3,
			"Description": fetchedproducts.Description,
			"Category":    fetchedproducts.Category.Name,
			"Price":       fetchedproducts.Price,
			"Size":        fetchedproducts.Size,
			"Quantity":    fetchedproducts.Quantity,
		}
		products = append(products, productdetails)
	}
	c.JSON(http.StatusOK, products)
}

func Addproduct(c *gin.Context) {
	err := c.BindJSON(&Product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}

	if Product.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Please enter a valid price"})
		return
	}

	if Product.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Please enter a valid quantity"})
		return
	}

	var dbproduct model.Product
	result := database.DB.Where("product_name=?", Product.Product_name).Find(&dbproduct)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch product name"})
		return
	}

	if dbproduct.Product_name == Product.Product_name {
		c.JSON(http.StatusConflict, gin.H{"Error": "This product already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Please upload image"})
}

func Upload(c *gin.Context) {
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to fetch images"})
		return
	}

	files := file.File["images"]
	var imagepaths []string

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "No images uploaded"})
		return
	}

	for _, val := range files {
		filepath := "./images/" + val.Filename
		if err := c.SaveUploadedFile(val, filepath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to upload image"})
			return
		}
		imagepaths = append(imagepaths, filepath)
	}

	Product.ImagePath1 = imagepaths[0]
	Product.ImagePath2 = imagepaths[1]
	Product.ImagePath3 = imagepaths[2]

	if len(imagepaths) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Please upload at least three images"})
		return
	}

	upload := database.DB.Create(&Product)
	if upload.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Product already exists"})
		return
	}
	Product = model.Product{}
	c.JSON(http.StatusOK, gin.H{"Message": "Product created successfully"})
}

func Blockuser(c *gin.Context) {
	var user model.UserModel
	id := c.Param("ID")
	result := database.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch id"})
		return
	}

	if user.Status {
		database.DB.Model(&user).Update("status", false)
		c.JSON(http.StatusOK, gin.H{"Message": "User Blocked"})
	} else {
		database.DB.Model(&user).Update("status", true)
		c.JSON(http.StatusOK, gin.H{"Error": "User Unblocked"})
	}
}

func Editcategory(c *gin.Context) {
	var edit model.Category
	id := c.Param("ID")
	result := c.BindJSON(&edit)
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}
	var editcategory model.Category
	fetch := database.DB.First(&editcategory, id)
	if fetch.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch category"})
		return
	}

	if fetch.RowsAffected > 0 {
		database.DB.Model(&editcategory).Updates(edit)
		c.JSON(http.StatusOK, gin.H{"Message": "Category updated successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update category"})
	}
}

func Editproduct(c *gin.Context) {
	var edit model.Product
	id := c.Param("ID")
	err := c.BindJSON(&edit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}
	var editproduct model.Product
	fetch := database.DB.First(&editproduct, id)
	if fetch.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch product"})
		return
	}

	if fetch.RowsAffected > 0 {
		database.DB.Model(&editproduct).Updates(edit)
		c.JSON(http.StatusOK, gin.H{"Message": "Product updated successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update product"})
	}
}

func Deletecategory(c *gin.Context) {
	var delete model.Category
	id := c.Param("ID")
	err := database.DB.First(&delete, id)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch"})
		return
	}
	fetch := database.DB.Model(&delete).Update("DeletedAt", time.Now())
	if fetch.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete category"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "Category deleted successfully"})
	}
}

func Deleteproduct(c *gin.Context) {
	var delete model.Product
	id := c.Param("ID")
	err := database.DB.First(&delete, id)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch"})
		return
	}
	fetch := database.DB.Model(&delete).Update("DeletedAt", time.Now())
	if fetch.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete product"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "Product deleted succcessfully"})
	}
}

func AdminLogout(c *gin.Context) {
	tokenstring := c.GetHeader("Authorization")
	if tokenstring == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Token not provided"})
		return
	}
	jwt.BlacklistedToken[tokenstring] = true
	c.JSON(http.StatusOK, gin.H{
		"Message":   "Admin logout successfully",
		"Blacklist": jwt.BlacklistedToken[tokenstring],
	})
}

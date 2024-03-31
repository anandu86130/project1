package routes

import (
	"fmt"
	"net/http"
	"project1/database"
	"project1/jwt"
	"project1/model"
	"project1/otp"
	"project1/send"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var Userdetails model.UserModel

const RoleUser = "user"

func Signup(c *gin.Context) {
	err := c.ShouldBindJSON(&Userdetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind"})
		return
	}

	var existinguser model.UserModel
	result := database.DB.Where("email=?", Userdetails.Email).First(&existinguser)
	if result.Error == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this user already exists"})
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(Userdetails.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hashpassword"})
		return
	}
	Userdetails.Password = string(hashedpassword)
	Userdetails.Status = true

	otp := otp.GenerateOTP(6)
	newOTP := model.OTP{
		Email: Userdetails.Email,
		Otp:   otp,
		Exp:   time.Now().Add(1 * time.Minute),
	}
	fmt.Println(otp)
	if err := database.DB.Where("email = ?", Userdetails.Email).First(&existinguser); err.Error != nil {
		database.DB.Model(&Userdetails).Updates(model.OTP{
			Otp: otp,
		})
	}
	if err := database.DB.Create(&newOTP).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate otp"})
		return
	}

	send.SendOTPByEmail(newOTP.Email, newOTP.Otp)
	c.JSON(http.StatusOK, gin.H{"message": "OTP send succcessfully"})
}

func Otpsignup(c *gin.Context) {
	var otp model.OTP
	err := c.BindJSON(&otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind"})
		return
	}

	var existingotp model.OTP
	result := database.DB.Where("email=?", Userdetails.Email).First(&existingotp)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch otp"})
		return
	}

	currentTime := time.Now()
	if currentTime.After(existingotp.Exp) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "otp expired"})
		return
	}

	if existingotp.Otp != otp.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid otp"})
		return
	}

	create := database.DB.Create(&Userdetails)
	fmt.Println(Userdetails)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
	}
}

func ResendOtp(c *gin.Context) {
	var fetch model.OTP
	err := c.BindJSON(&fetch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch email"})
		return
	}

	var existinguser model.OTP
	fetcheddata := database.DB.Where("email=?", fetch.Email).First(&existinguser)
	if fetcheddata.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email not found"})
		return
	}

	newOTP := otp.GenerateOTP(6)

	result := database.DB.Model(&model.OTP{}).Where("email=?", fetch.Email).Updates(model.OTP{Otp: newOTP, Exp: time.Now().Add(1 * time.Minute)})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resend otp"})
		return
	}

	send.SendOTPByEmail(fetch.Email, newOTP)

	c.JSON(http.StatusOK, gin.H{"message": "OTP resent successfully"})
}

func Login(c *gin.Context) {
	var userlogin model.UserModel
	err := c.ShouldBindJSON(&userlogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind"})
		return
	}

	var existinguser model.UserModel
	result := database.DB.Where("email=?", userlogin.Email).First(&existinguser)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect email or password"})
		return
	}

	password := bcrypt.CompareHashAndPassword([]byte(existinguser.Password), []byte(userlogin.Password))
	if password != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	} else {
		if existinguser.Status {
			jwt.JwtToken(c, existinguser.UserID, userlogin.Email, RoleUser)
			c.JSON(http.StatusOK, gin.H{"message": "Login successfully"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "blocked user"})
		}
	}

}

func Productview(c *gin.Context) {
	var product []model.Product
	result := database.DB.Find(&product)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to load"})
		return
	}
	var productview []gin.H
	for _, fetchedproducts := range product {
		details := gin.H{
			"productid":   fetchedproducts.ID,
			"name":        fetchedproducts.Product_name,
			"imagepath1":  fetchedproducts.ImagePath1,
			"imagepath2":  fetchedproducts.ImagePath2,
			"imagepath3":  fetchedproducts.ImagePath3,
			"description": fetchedproducts.Description,
			"price":       fetchedproducts.Price,
			"size":        fetchedproducts.Size,
			"quantity":    fetchedproducts.Quantity,
		}
		productview = append(productview, details)
	}
	c.JSON(http.StatusOK, productview)
}

func Productdetails(c *gin.Context) {
	id := c.Param("ID")
	productid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product id"})
		return
	}
	var Product model.Product
	if err := database.DB.Where("id=?", productid).First(&Product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}
	type productdetails struct {
		Name        string
		Price       uint
		Imagepath1  string
		ImagePath2  string
		ImagePath3  string
		Description string
		Size        string
	}

	details := productdetails{
		Name:        Product.Product_name,
		Price:       Product.Price,
		Imagepath1:  Product.ImagePath1,
		ImagePath2:  Product.ImagePath2,
		ImagePath3:  Product.ImagePath3,
		Description: Product.Description,
		Size:        Product.Size,
	}

	var rating model.Rating
	database.DB.Where("product_id=?", productid).Find(&rating)

	type ratingdetails struct {
		Rating uint
		Review string
	}

	rdetails := ratingdetails{
		Rating: rating.Rating,
		Review: rating.Review,
	}
	c.JSON(http.StatusOK, gin.H{
		"Product details":           details,
		"product rating and review": rdetails,
	})
}

func Productsearch(c *gin.Context) {
	search := c.Request.FormValue("search")
	if search == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Please send the search"})
		return
	}
	var product []model.Product
	if result := database.DB.Where("product_name ILIKE ?", "%"+search+"%").Find(&product).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find product"})
		return
	}
	c.JSON(http.StatusOK, product)
}
func Logout(c *gin.Context) {
	tokenstring := c.GetHeader("Authorization")
	if tokenstring == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Token not found"})
		return
	}

	jwt.Userdetails = model.UserModel{}
	jwt.BlacklistedToken[tokenstring] = true

	c.JSON(http.StatusOK, gin.H{
		"Message":   "Successfully logout",
		"Blacklist": jwt.BlacklistedToken[tokenstring],
	})
}

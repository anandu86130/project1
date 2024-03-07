package routes

import (
	"fmt"
	"net/http"
	"project1/database"
	"project1/model"
	"project1/otp"
	"project1/send"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var Userdetails model.UserModel

func Signup(c *gin.Context) {
	err := c.ShouldBindJSON(&Userdetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	var existinguser model.UserModel
	result := database.DB.Where("email=?", Userdetails.Email).First(&existinguser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, "this user already exists")
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(Userdetails.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to hashpassword")
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
	if err := database.DB.Where("email = ?", Userdetails.Email).First(&existinguser); err.Error == nil {
		database.DB.Model(&Userdetails).Updates(model.OTP{
			Otp: otp,
		})
	} else {
		if err := database.DB.Create(&newOTP).Error; err != nil {
			c.JSON(http.StatusInternalServerError, "failed to generate otp")
			return
		}
	}

	send.SendOTPByEmail(newOTP.Email, newOTP.Otp)
	c.JSON(http.StatusOK, "OTP send successfully")
}

func Otpsignup(c *gin.Context) {
	var otp model.OTP
	err := c.BindJSON(&otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	var existingotp model.OTP
	result := database.DB.Where("otp=?", otp.Otp).First(&existingotp)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to fetch otp")
		return
	}

	currentTime := time.Now()
	if currentTime.After(existingotp.Exp) {
		c.JSON(http.StatusInternalServerError, "otp expired")
		return
	}

	if existingotp.Otp != otp.Otp {
		c.JSON(http.StatusBadRequest, "invalid otp")
		return
	} else {
		create := database.DB.Create(&Userdetails)
		fmt.Println(Userdetails)
		if create.Error != nil {
			c.JSON(http.StatusInternalServerError, "failed to create user")
			return
		}
		c.JSON(http.StatusOK, "user created successfully")
	}
}

func ResendOtp(c *gin.Context) {
	var fetch model.OTP
	err := c.BindJSON(&fetch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to fetch user")
		return
	}

	var existinguser model.OTP
	fetcheddata := database.DB.Where("email=?", fetch.Email).First(&existinguser)
	if fetcheddata.Error != nil {
		c.JSON(http.StatusBadRequest, "user not found")
		return
	}

	newOTP := otp.GenerateOTP(6)

	result := database.DB.Model(&model.OTP{}).Where("email=?", fetch.Email).Updates(model.OTP{Otp: newOTP, Exp: time.Now().Add(1 * time.Minute)})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to update otp")
		return
	}

	send.SendOTPByEmail(fetch.Email, newOTP)

	c.JSON(http.StatusOK, "OTP resent successfully")
}

func Login(c *gin.Context) {
	var userlogin model.UserModel
	err := c.ShouldBindJSON(&userlogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	var existinguser model.UserModel
	email := database.DB.Where("email=?", userlogin.Email).First(&existinguser)
	if email.Error != nil {
		c.JSON(http.StatusUnauthorized, "incorrect email or password")
		return
	}

	result := bcrypt.CompareHashAndPassword([]byte(existinguser.Password), []byte(userlogin.Password))
	if result != nil {
		c.JSON(http.StatusUnauthorized, "invalid email or password")
		return
	} else {
		if existinguser.Status {
			c.JSON(http.StatusOK, "Login successfully")
		} else {
			c.JSON(http.StatusUnauthorized, "blocked user")
		}
	}

}

func Productview(c *gin.Context) {
	var product []model.Product
	result := database.DB.Find(&product)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, "failed to load")
		return
	}
	var productview []gin.H
	for _, fetchedproducts := range product {
		details := gin.H{
			"id":         fetchedproducts.ProductId,
			"name":       fetchedproducts.Product_name,
			"imagepath1": fetchedproducts.ImagePath1,
			"imagepath2": fetchedproducts.ImagePath2,
			"imagepath3": fetchedproducts.ImagePath3,
			"price":      fetchedproducts.Price,
			"size":       fetchedproducts.Size,
			"quantity":   fetchedproducts.Quantity,
		}
		productview = append(productview, details)
	}
	c.JSON(http.StatusOK, productview)
}

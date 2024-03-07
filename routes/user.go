package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var userdetails model.UserModel
	err := c.ShouldBindJSON(&userdetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(userdetails.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to hashpassword")
		return
	}
	userdetails.Password = string(hashedpassword)
	userdetails.Status = true

	result := database.DB.Create(&userdetails)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, "user already exists")
		return
	}

	c.JSON(http.StatusOK, "successfully signed up")
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
	}else{
		if existinguser.Status{
			c.JSON(http.StatusOK, "Login successfully")
		}else{
			c.JSON(http.StatusUnauthorized, "invalid email or password")
		}
	}

	
}

// func Otp(c *gin.Context) {
// 	var otp model.OTP
// 	err := c.ShouldBindJSON(&otp)
// 	if err != nil {
// 		c.JSON(500, "failed to bind")
// 	}
// }
// func Verifyotp(c *gin.Context) {

// }

func Productview(c *gin.Context) {

}

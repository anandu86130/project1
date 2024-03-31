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

var UserCheck model.UserModel

func Forgotpassword(c *gin.Context) {
	var Otpstore model.OTP
	err := c.ShouldBindJSON(&UserCheck)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}

	if err := database.DB.First(&UserCheck, "email=?", UserCheck.Email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "User not found"})
		return
	}
	Otp := otp.GenerateOTP(6)
	Otpstore.Otp = Otp
	Otpstore.Email = UserCheck.Email
	Otpstore.Exp = time.Now().Add(1 * time.Minute)
	fmt.Println("----------------", Otp, "-----------------")

	send.SendOTPByEmail(UserCheck.Email, Otp)

	result := database.DB.First(&Otpstore, "email=?", UserCheck.Email)
	if result.Error != nil {
		Otpstore = model.OTP{
			Otp:   Otp,
			Email: UserCheck.Email,
			Exp:   time.Now().Add(30 * time.Second),
		}
		err := database.DB.Create(&Otpstore)
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save otp details"})
		}
	} else {
		err := database.DB.Model(&Otpstore).Where("email=?", UserCheck.Email).Updates(model.OTP{
			Otp: Otp,
			Exp: time.Now().Add(1 * time.Minute),
		})
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update data"})
		}
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Otp send to email"})
}

func Otpcheck(c *gin.Context) {
	var otp model.OTP
	err := c.ShouldBindJSON(&otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}
	var checkotp model.OTP
	result := database.DB.Where("email=?", UserCheck.Email).First(&checkotp)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find email"})
		return
	}

	checking := database.DB.Where("otp=?", otp.Otp).Find(&checkotp)
	if checking.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find otp"})
		return
	}

	if otp.Otp != checkotp.Otp {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Invalid otp"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "OTP is correct, you can now change the password"})
}

func PasswordReset(c *gin.Context) {
	var password model.UserModel
	err := c.ShouldBindJSON(&password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "failed to hashpassword"})
		return
	}
	UserCheck.Password = string(hashedpassword)

	result := database.DB.Model(&password).Where("email=?", UserCheck.Email).Updates(model.UserModel{
		Password: UserCheck.Password,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Password changed successfully"})
}

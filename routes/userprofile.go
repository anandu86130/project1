package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

type Person struct {
	UserID uint
	Name   string
	Email  string
}

type UserAddress struct {
	AddressId uint
	Address   string
	City      string
	Landmark  string
	State     string
	Country   string
	Pincode   *string
}

func UserProfile(c *gin.Context) {
	var user model.UserModel
	userid := c.GetUint("userid")
	result := database.DB.Where("user_id=?", userid).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find user"})
		return
	}
	var address []model.Address
	var Addressshow []UserAddress
	err := database.DB.Where("user_id", userid).Find(&address)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find user"})
		return
	}

	users := Person{
		UserID: userid,
		Name:   user.Name,
		Email:  user.Email,
	}

	for _, A := range address {
		Addressshow = append(Addressshow, UserAddress{A.AddressId, A.Address, A.City, A.Landmark, A.State, A.Country, &A.Pincode})
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    users,
		"address": Addressshow,
	})
}

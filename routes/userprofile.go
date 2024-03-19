package routes

import (
	"fmt"
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name  string
	Email string
}

type UserAddress struct {
	Address  string
	City     string
	Landmark string
	State    string
	Country  string
	Pincode  *string
}

func UserProfile(c *gin.Context) {
	var user model.UserModel
	userid := c.GetUint("userid")
	result := database.DB.Where("user_id=?", userid).First(&user)
	fmt.Println("user========================================================================================", userid)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}
	var address []model.Address
	var Addressshow []UserAddress
	err := database.DB.Where("user_id", userid).First(&address)
	fmt.Println("addresssssssssssss================================================================", userid)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	users := Person{
		Name:  user.Name,
		Email: user.Email,
	}

	for _, A := range address {
		Addressshow = append(Addressshow, UserAddress{A.Address, A.City, A.Landmark, A.State, A.Country, &A.Pincode})
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"user":    users,
		"address": Addressshow,
	})
}

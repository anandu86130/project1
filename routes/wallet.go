package routes

import (
	"fmt"
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func Wallet(c *gin.Context) {
	userid := c.GetUint("userid")
	fmt.Println("=========", userid)
	var wallet model.Wallet
	if err := database.DB.Where("user_id=?", userid).First(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Wallet amount": wallet.Amount,
	})
}

package invoice

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func Invoicedownload(c *gin.Context) {
	userid := c.GetUint("userid")
	var user model.UserModel
	if err := database.DB.Where("user_id=?", userid).First(&user).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"Error":"Failed to find user"})
		return
	}
}

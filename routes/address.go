package routes

import (
	"net/http"
	"project1/database"
	"project1/model"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {
	id := c.Param("ID")
	var user model.UserModel
	result := database.DB.Where("user_id=?", id).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find user"})
		return
	}

	var address model.Address
	err := c.BindJSON(&address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}

	var dbaddress model.Address
	database.DB.Where("address=?", address.Address).First(&dbaddress)

	if dbaddress.AddressId != 0 && address.Address == dbaddress.Address {
		c.JSON(http.StatusConflict, gin.H{"Error": "This address already exists"})
		return
	}

	if len(address.Pincode) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid pincode"})
		return
	}

	create := database.DB.Create(&model.Address{
		Address:  address.Address,
		City:     address.City,
		Landmark: address.Landmark,
		State:    address.State,
		Country:  address.Country,
		Pincode:  address.Pincode,
		UserId:   user.UserID,
	})
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Address added successfully"})
}

func EditAddress(c *gin.Context) {
	id := c.Param("ID")
	var address model.Address
	update := c.BindJSON(&address)
	if update != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to bind"})
		return
	}
	var dbaddress model.Address
	result := database.DB.First(&dbaddress, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Address not found"})
		return
	}

	if result.RowsAffected > 0 {
		database.DB.Model(&dbaddress).Updates(address)
		c.JSON(http.StatusOK, gin.H{"Message": "Address updated successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update address"})
		return
	}
}

func Deleteaddress(c *gin.Context) {
	addressID := c.Param("ID")
	var dbaddress model.Address
	result := database.DB.Where("address_id=?", addressID).First(&dbaddress)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Address not found"})
		return
	}
	delete := database.DB.Delete(&dbaddress)
	if delete.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Address deleted successfully"})
}

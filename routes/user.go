package routes

import (
	"fmt"
	"net/http"
	"project1/database"
	"project1/jwt"
	"project1/model"
	"project1/otp"
	"project1/send"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var Userdetails model.UserModel
const RoleUser = "user"

func Signup(c *gin.Context) {
	err := c.ShouldBindJSON(&Userdetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	var existinguser model.UserModel
	result := database.DB.Where("email=?", Userdetails.Email).First(&existinguser)
	if result.Error == nil {
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
	c.JSON(http.StatusOK, "OTP send succcessfully")
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
	}

	create := database.DB.Create(&Userdetails)
	fmt.Println(Userdetails)
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to create user")
		return
	} else {
		c.JSON(http.StatusOK, "user created successfully")
	}

}

// func Profile(c *gin.Context) {
// 	var user []model.UserModel
// 	result := database.DB.Find(&user)
// 	if result.Error != nil {
// 		c.JSON(http.StatusBadRequest, "failed to find user")
// 		return
// 	}

// 	var users []gin.H
// 	for _, details := range user {
// 		userdata := gin.H{
// 			"name": details.Name,
// 		}
// 		users = append(users, userdata)
// 	}
// 	c.JSON(http.StatusOK, users)
// }

func AddAddress(c *gin.Context) {
	// user, _ := c.Get("ID")
	// userID := user.(model.UserModel).ID

	var address model.Address
	err := c.BindJSON(&address)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	var dbaddress model.Address
	// dbaddress.User_ID = userID
	database.DB.Where("address=?", address.Address).First(&dbaddress)

	// if dbaddress.ID != 0 && address.Address == dbaddress.Address {
	// 	c.JSON(http.StatusConflict, "this address already exists")
	// 	return
	// }

	if len(address.Pincode) != 6 {
		c.JSON(http.StatusInternalServerError, "invalid pincode")
		return
	}

	create := database.DB.Create(&model.Address{
		Address:  address.Address,
		City:     address.City,
		Landmark: address.Landmark,
		State:    address.State,
		Country:  address.Country,
		Pincode:  address.Pincode,
		// User_ID:  address.User_ID,
	})
	if create.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to create address")
		return
	}

	c.JSON(http.StatusOK, "Address added successfully")
}

func EditAddress(c *gin.Context) {
	var address model.Address
	err := c.BindJSON(&address)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to bind")
		return
	}

	addressID := c.Param("ID")

	var dbaddress model.Address
	result := database.DB.Where("id=?", addressID).First(&dbaddress)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, "address not found")
		return
	}

	dbaddress.Address = address.Address
	dbaddress.City = address.City
	dbaddress.Country = address.City
	dbaddress.Landmark = address.Landmark
	dbaddress.Pincode = address.Pincode
	dbaddress.State = address.State

	update := database.DB.Save(&dbaddress)
	if update.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to update address")
		return
	}

	c.JSON(http.StatusOK, "address updated successfully")
}

func Deleteaddress(c *gin.Context) {
	addressID := c.Param("ID")
	var dbaddress model.Address
	result := database.DB.Where("id=?", addressID).First(&dbaddress)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, "address not found")
		return
	}
	delete := database.DB.Delete(&dbaddress)
	if delete.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to delete address")
		return
	}

	c.JSON(http.StatusOK, "address deleted successfully")
}

func ResendOtp(c *gin.Context) {
	var fetch model.OTP
	err := c.BindJSON(&fetch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to fetch email")
		return
	}

	var existinguser model.OTP
	fetcheddata := database.DB.Where("email=?", fetch.Email).First(&existinguser)
	if fetcheddata.Error != nil {
		c.JSON(http.StatusBadRequest, "email not found")
		return
	}

	newOTP := otp.GenerateOTP(6)

	result := database.DB.Model(&model.OTP{}).Where("email=?", fetch.Email).Updates(model.OTP{Otp: newOTP, Exp: time.Now().Add(1 * time.Minute)})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, "failed to resend otp")
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
			jwt.JwtToken(c, userlogin.ID, userlogin.Email, RoleUser)
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
			"id":          fetchedproducts.ProductId,
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

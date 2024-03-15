package main

import (
	"project1/database"
	"project1/helper"
	"project1/jwt"
	"project1/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	helper.LoadEnv()
	database.DBconnect()
}

func main() {
	r := gin.Default()

	//user
	r.POST("/user/signup", routes.Signup)
	r.POST("/user/otpsignup", routes.Otpsignup)
	r.POST("/user/resendotp", routes.ResendOtp)
	r.POST("/user/login", routes.Login)
	r.POST("user/forgotpassword",jwt.AuthMiddleware("user"), routes.Forgotpassword)
	r.POST("/user/forgototpcheck",jwt.AuthMiddleware("user"), routes.Otpcheck)
	r.POST("/user/resetpassword", jwt.AuthMiddleware("user"),routes.PasswordReset)
	// r.GET("/user/profile", routes.Profile)
	r.POST("/user/address", jwt.AuthMiddleware("user"), routes.AddAddress)
	r.PATCH("/user/address/:ID", jwt.AuthMiddleware("user"), routes.EditAddress)
	r.DELETE("/user/address/:ID", jwt.AuthMiddleware("user"), routes.Deleteaddress)
	r.GET("/user/product", jwt.AuthMiddleware("user"), routes.Productview)
	r.GET("/user/cart", jwt.AuthMiddleware("user"),routes.CartView)

	//admin
	r.POST("/admin/signin", routes.Signin)
	r.GET("/admin/getuser", jwt.AuthMiddleware("admin"), routes.Getuser)
	r.PATCH("/admin/blockuser/:ID", jwt.AuthMiddleware("admin"), routes.Blockuser)
	//category
	r.GET("/admin/category", jwt.AuthMiddleware("admin"), routes.Category)
	r.POST("/admin/category", jwt.AuthMiddleware("admin"), routes.Addcategory)
	r.PATCH("/admin/category/:ID", jwt.AuthMiddleware("admin"), routes.Editcategory)
	r.DELETE("/admin/category/:ID", jwt.AuthMiddleware("admin"), routes.Deletecategory)
	//product
	r.GET("/admin/product", jwt.AuthMiddleware("admin"), routes.Aproduct)
	r.POST("/admin/product", jwt.AuthMiddleware("admin"), routes.Addproduct)
	r.POST("/admin/upload", jwt.AuthMiddleware("admin"), routes.Upload)
	r.PATCH("/admin/product/:ID", jwt.AuthMiddleware("admin"), routes.Editproduct)
	r.DELETE("/admin/product/:ID", jwt.AuthMiddleware("admin"), routes.Deleteproduct)

	r.Run(":8080")
}

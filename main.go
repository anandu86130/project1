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

	//user authentication
	r.POST("/user/signup", routes.Signup)
	r.POST("/user/otpsignup", routes.Otpsignup)
	r.POST("/user/resendotp", routes.ResendOtp)
	r.POST("/user/login", routes.Login)
	//user profile
	r.GET("/user/profile", jwt.AuthMiddleware("user"), routes.UserProfile)
	//user forgotpassword
	r.POST("user/forgotpassword", routes.Forgotpassword)
	r.POST("/user/forgototpcheck", routes.Otpcheck)
	r.POST("/user/resetpassword", routes.PasswordReset)
	//user address
	r.POST("/user/address/:ID", jwt.AuthMiddleware("user"), routes.AddAddress)
	r.PATCH("/user/address/:ID", jwt.AuthMiddleware("user"), routes.EditAddress)
	r.DELETE("/user/address/:ID", jwt.AuthMiddleware("user"), routes.Deleteaddress)
	//user product
	r.GET("/user/product", jwt.AuthMiddleware("user"), routes.Productview)
	//user cart
	r.GET("/user/cart", jwt.AuthMiddleware("user"), routes.CartView)
	r.POST("/user/cart/:ID", jwt.AuthMiddleware("user"), routes.Addtocart)
	r.DELETE("/user/cart/:ID", jwt.AuthMiddleware("user"), routes.Deletecart)
	//user checkout
	r.POST("user/checkout/:ID", jwt.AuthMiddleware("user"), routes.Checkout)
	//user order
	r.GET("/user/order", jwt.AuthMiddleware("user"), routes.Orderview)
	r.POST("user/order/:ID", jwt.AuthMiddleware("user"), routes.Orderdetails)
	// r.DELETE("user/order", jwt.AuthMiddleware("user"), routes.Cancelorder)
	//user logout
	r.GET("/user/logout", jwt.AuthMiddleware("user"), routes.Logout)

	//admin authentucation
	r.POST("/admin/signin", routes.Signin)
	// admin user management
	r.GET("/admin/getuser", jwt.AuthMiddleware("admin"), routes.Getuser)
	r.PATCH("/admin/blockuser/:ID", jwt.AuthMiddleware("admin"), routes.Blockuser)
	//admin category management
	r.GET("/admin/category", jwt.AuthMiddleware("admin"), routes.Category)
	r.POST("/admin/category", jwt.AuthMiddleware("admin"), routes.Addcategory)
	r.PATCH("/admin/category/:ID", jwt.AuthMiddleware("admin"), routes.Editcategory)
	r.DELETE("/admin/category/:ID", jwt.AuthMiddleware("admin"), routes.Deletecategory)
	//admin product mangement
	r.GET("/admin/product", jwt.AuthMiddleware("admin"), routes.Aproduct)
	r.POST("/admin/product", jwt.AuthMiddleware("admin"), routes.Addproduct)
	r.POST("/admin/upload", jwt.AuthMiddleware("admin"), routes.Upload)
	r.PATCH("/admin/product/:ID", jwt.AuthMiddleware("admin"), routes.Editproduct)
	r.DELETE("/admin/product/:ID", jwt.AuthMiddleware("admin"), routes.Deleteproduct)
	//admin coupon
	r.GET("/admin/coupon", jwt.AuthMiddleware("admin"), routes.Coupon)
	r.POST("admin/coupon", jwt.AuthMiddleware("admin"), routes.Addcoupon)
	r.DELETE("admin/coupon/:ID", jwt.AuthMiddleware("admin"), routes.Deletecoupon)
	//admin logout
	r.GET("admin/logout", jwt.AuthMiddleware("admin"), routes.AdminLogout)

	r.Run(":8080")
}

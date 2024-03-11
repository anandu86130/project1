package main

import (
	"project1/database"
	"project1/helper"
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
	r.GET("/user/profile", routes.Profile)
	r.GET("/user/product", routes.Productview)

	//admin
	r.POST("/admin/signin", routes.Signin)
	r.GET("/admin/getuser", routes.Getuser)
	r.PATCH("/admin/blockuser/:ID", routes.Blockuser)
	//category
	r.GET("/admin/category", routes.Category)
	r.POST("/admin/category", routes.Addcategory)
	r.PATCH("/admin/category/:ID", routes.Editcategory)
	r.DELETE("/admin/category/:ID", routes.Deletecategory)
	//product
	r.GET("/admin/product", routes.Aproduct)
	r.POST("/admin/product", routes.Addproduct)
	r.POST("/admin/upload", routes.Upload)
	r.PATCH("/admin/product/:ID", routes.Editproduct)
	r.DELETE("/admin/product/:ID", routes.Deleteproduct)

	r.Run(":8080")
}

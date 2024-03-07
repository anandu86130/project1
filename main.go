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
	r.POST("/signup", routes.Signup)
	r.POST("/login", routes.Login)
	r.GET("/productview", routes.Productview)
	// r.POST("/otp", routes.Otp)

	//admin
	r.POST("/signin", routes.Signin)
	r.GET("/getuser", routes.Getuser)
	r.PATCH("/blockuser/:ID", routes.Blockuser)
	//category
	// r.GET("/category", routes.Category)
	r.GET("/category", routes.Category)
	r.POST("/addcategory", routes.Addcategory)
	r.PUT("/editcategory", routes.Editcategory)
	r.DELETE("/deletecategory", routes.Deletecategory)
	//product
	// r.GET("/aproduct", routes.Aproduct)
	r.GET("/product", routes.Aproduct)
	r.POST("/addproduct", routes.Addproduct)
	r.PUT("/editproduct", routes.Editproduct)
	r.DELETE("/deleteproduct", routes.Deleteproduct)

	r.Run(":8080")
}

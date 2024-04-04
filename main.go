package main

import (
	"net/http"
	"project1/database"
	"project1/helper"
	"project1/invoice"
	"project1/jwt"
	"project1/payment"
	"project1/routes"
	"project1/sales"

	"github.com/gin-gonic/gin"
)

func init() {
	helper.LoadEnv()
	database.DBconnect()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

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
	r.POST("/user/product/:ID", jwt.AuthMiddleware("user"), routes.Productdetails)
	//user productsearch
	r.POST("/user/search", jwt.AuthMiddleware("user"), routes.Productsearch)
	//user cart
	r.GET("/user/cart", jwt.AuthMiddleware("user"), routes.CartView)
	r.POST("/user/cart/:ID", jwt.AuthMiddleware("user"), routes.Addtocart)
	r.DELETE("/user/cart/:ID", jwt.AuthMiddleware("user"), routes.Deletecart)
	//user checkout
	r.POST("user/checkout/:ID", jwt.AuthMiddleware("user"), routes.Checkout)
	//user order
	r.GET("/user/order", jwt.AuthMiddleware("user"), routes.Orderview)
	r.POST("/user/order/:ID", jwt.AuthMiddleware("user"), routes.Orderdetails)
	r.PATCH("/user/order/:ID", jwt.AuthMiddleware("user"), routes.Cancelorder)
	//user payment
	r.GET("/user/payment", func(c *gin.Context) {
		c.HTML(http.StatusOK, "payment.html", nil)
	})
	r.POST("/user/payment/confirm", payment.Paymentconfirmation)
	//user whishlist
	r.GET("/user/whishlist", jwt.AuthMiddleware("user"), routes.Whishlist)
	r.POST("/user/whishlist/:ID", jwt.AuthMiddleware("user"), routes.Addtowhishlist)
	r.DELETE("/user/whishlist/:ID", jwt.AuthMiddleware("user"), routes.Deletewhishlist)
	//user rating
	r.POST("/user/rating/:ID", jwt.AuthMiddleware("user"), routes.Addrating)
	//user product sort
	r.POST("/user/product/sort", jwt.AuthMiddleware("user"), routes.Sortproduct)
	//user category filter
	r.POST("/user/product/category", jwt.AuthMiddleware("user"), routes.Categoryfilter)
	//user wallet
	r.GET("/user/wallet", jwt.AuthMiddleware("user"), routes.Wallet)
	// admin invoice
	r.GET("/user/invoice/:ID", jwt.AuthMiddleware("user"), invoice.Invoicedownload)
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
	r.POST("/admin/coupon", jwt.AuthMiddleware("admin"), routes.Addcoupon)
	r.DELETE("/admin/coupon/:ID", jwt.AuthMiddleware("admin"), routes.Deletecoupon)
	//admin order management
	r.GET("/admin/order", jwt.AuthMiddleware("admin"), routes.Adminorderview)
	r.PATCH("/admin/order/:ID", jwt.AuthMiddleware("admin"), routes.Adminorderstatus)
	r.PUT("/admin/order/:ID", jwt.AuthMiddleware("admin"), routes.Admincancelorders)
	//admin addproductoffer
	r.POST("/admin/productoffer", jwt.AuthMiddleware("admin"), routes.Addproductoffer)
	//admin salesreport
	r.GET("/admin/salesreport", jwt.AuthMiddleware("admin"), sales.Salesreport)
	//admin logout
	r.GET("/admin/logout", jwt.AuthMiddleware("admin"), routes.AdminLogout)

	r.Run(":8080")
}

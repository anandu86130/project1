package jwt

import (
	"fmt"
	"net/http"
	"project1/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SecretKey = []byte("qwertuiouplkhgfdsazxcvbnm")
var Userdetails model.UserModel
var BlacklistedToken = make(map[string]bool)

type Claims struct {
	ID    uint
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// func JwtTokenStart(c *gin.Context, userId uint, email string, role string) {
// 	tokenString, err := createToken(userId, email, role)
// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"Error": "Failed to create Token",
// 		})
// 	}
// 	c.Set("token", tokenString)
// 	c.JSON(201, gin.H{
// 		"Token": tokenString,
// 	})
// 	fmt.Println("---------------===  ", tokenString, "  ===-----------------")
// }

func JwtToken(c *gin.Context, id uint, email string, role string) {
	claims := Claims{
		ID:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring := c.GetHeader("Authorization")
		if tokenstring == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}
		if BlacklistedToken[tokenstring] {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token removed"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No permission"})
			c.Abort()
			return
		}

		c.Set("userid", claims.ID)
		fmt.Println("userid============================================================================================", claims.ID)
		c.Next()
	}
}

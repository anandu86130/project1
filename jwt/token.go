package jwt

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SecretKey = []byte("qwertuiouplkhgfdsazxcvbnm")
var BlacklistedTokens = make(map[string]bool)

type Claims struct {
	ID    uint   
	Email string `json:"email"`
	Role  string `json:"role"`

	jwt.StandardClaims
}

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
	if BlacklistedTokens[tokenstring] {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token revoked"})
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
	c.Next()
}
}

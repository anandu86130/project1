package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secretqlusilomasqwi")

func CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtKey)
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
	if err != nil {
		return nil, err
	}
	return token, nil
}

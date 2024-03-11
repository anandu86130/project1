package otp

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateOTP(length int) string {
	characters := "0123456789"
	otp := ""
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(characters))
		otp += string(characters[randomIndex])
	}

	return otp
}

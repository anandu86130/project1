package helper

import (
	"fmt"

	"github.com/joho/godotenv"
)

// load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load env")
	}
}

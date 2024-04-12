package helper

import "github.com/joho/godotenv"

//load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env")
	}
}

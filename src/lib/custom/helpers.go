package custom

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	log.Println("Loading Environment variables from .env")
	return godotenv.Load("../.env")
}

// GetEnv returns the value of the environment 
// variable tagged to the requested key.
func GetEnv(key string) string {
	return os.Getenv(key)
}

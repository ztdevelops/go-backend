package custom

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
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

// GetOpt returns the Client options that are derived
// from their respective JSON files.
func GetOpt(v string) option.ClientOption {
	return option.WithCredentialsFile(GetEnv(v))
}
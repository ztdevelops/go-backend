package middleware

import (
	"github.com/rs/cors"
)

func GetCorsWrapper(allowedHeaders, allowedMethods []string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedMethods: allowedMethods,
		AllowedHeaders: allowedHeaders,
	})
}
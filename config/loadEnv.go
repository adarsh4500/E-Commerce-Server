package config

import (
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret string
var DB_UserName string
var DB_Password string
var DB_Name string
var DB_Host string
var DB_Port string
var DB_SSLMode string
var AllowedOrigins string

func LoadEnv() error {
	_ = godotenv.Load()

	JWTSecret = os.Getenv("JWT_SECRET_KEY")
	DB_UserName = os.Getenv("DB_USERNAME")
	DB_Password = os.Getenv("DB_PASSWORD")
	DB_Name = os.Getenv("DB_NAME")
	DB_Host = os.Getenv("DB_HOST")
	DB_Port = os.Getenv("DB_PORT")
	DB_SSLMode = os.Getenv("DB_SSLMODE")
	AllowedOrigins = os.Getenv("ALLOWED_ORIGINS")

	return nil
}

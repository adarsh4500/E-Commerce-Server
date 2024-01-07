package config

import (
	"os"

	"github.com/joho/godotenv"
)

var JWTSecret string
var DB_UserName string
var DB_Password string
var DB_Name string

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	JWTSecret = os.Getenv("JWT_SECRET_KEY")
	DB_UserName = os.Getenv("DB_USERNAME")
	DB_Password = os.Getenv("DB_PASSWORD")
	DB_Name = os.Getenv("DB_NAME")

	return nil
}
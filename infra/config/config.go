package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type appConfig struct {
	DB_DIALECT  string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	APP_PORT    string
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("failed to load env file, err:", err.Error())
	}
}

func GetAppConfig() *appConfig {
	return &appConfig{
		DB_DIALECT:  os.Getenv("DB_DIALECT"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		APP_PORT:    os.Getenv("APP_PORT"),
	}
}

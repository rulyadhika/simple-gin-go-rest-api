package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type appConfig struct {
	DB_DIALECT                        string
	DB_HOST                           string
	DB_PORT                           string
	DB_NAME                           string
	DB_USER                           string
	DB_PASSWORD                       string
	APP_PORT                          string
	ACCESS_TOKEN_SECRET               string
	REFRESH_TOKEN_SECRET              string
	ACCESS_TOKEN_EXPIRATION_DURATION  time.Duration
	REFRESH_TOKEN_EXPIRATION_DURATION time.Duration
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("failed to load env file, err:", err.Error())
	}
}

func GetAppConfig() *appConfig {
	accessTokenExpDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRATION_DURATION"))

	if err != nil {
		log.Fatal("failed to parse string to time.Duration, err:", err.Error())
	}

	refreshTokenExpDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRATION_DURATION"))

	if err != nil {
		log.Fatal("failed to parse string to time.Duration, err:", err.Error())
	}

	return &appConfig{
		DB_DIALECT:                        os.Getenv("DB_DIALECT"),
		DB_HOST:                           os.Getenv("DB_HOST"),
		DB_PORT:                           os.Getenv("DB_PORT"),
		DB_NAME:                           os.Getenv("DB_NAME"),
		DB_USER:                           os.Getenv("DB_USER"),
		DB_PASSWORD:                       os.Getenv("DB_PASSWORD"),
		APP_PORT:                          os.Getenv("APP_PORT"),
		ACCESS_TOKEN_SECRET:               os.Getenv("ACCESS_TOKEN_SECRET"),
		REFRESH_TOKEN_SECRET:              os.Getenv("REFRESH_TOKEN_SECRET"),
		ACCESS_TOKEN_EXPIRATION_DURATION:  accessTokenExpDuration,
		REFRESH_TOKEN_EXPIRATION_DURATION: refreshTokenExpDuration,
	}
}

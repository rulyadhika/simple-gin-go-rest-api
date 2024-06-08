package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type appConfig struct {
	DB_DIALECT                                   string
	DB_HOST                                      string
	DB_PORT                                      uint
	DB_NAME                                      string
	DB_USER                                      string
	DB_PASSWORD                                  string
	APP_PORT                                     uint
	ACCESS_TOKEN_SECRET                          string
	REFRESH_TOKEN_SECRET                         string
	ACCESS_TOKEN_EXPIRATION_DURATION             time.Duration
	REFRESH_TOKEN_EXPIRATION_DURATION            time.Duration
	ACCOUNT_ACTIVATION_TOKEN_EXPIRATION_DURATION time.Duration
	ACCOUNT_ACTIVATION_URL                       string
	EMAIL_SMTP_SERVER                            string
	EMAIL_SMTP_PORT                              uint
	EMAIL_SENDER_IDENTITY                        string
	EMAIL_SENDER_AND_SMTP_USER                   string
	EMAIL_SMTP_USER_PASSWORD                     string
	PASSWORD_RESET_TOKEN_EXPIRATION_DURATION     time.Duration
	PASSWORD_RESET_TOKEN_LENGTH                  uint
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

	accountActivationTokenExpDuration, err := time.ParseDuration(os.Getenv("ACCOUNT_ACTIVATION_TOKEN_EXPIRATION_DURATION"))
	if err != nil {
		log.Fatal("failed to parse string to time.Duration, err:", err.Error())
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("failed to parse string to uint, err:", err.Error())
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatal("failed to parse string to uint, err:", err.Error())
	}

	emailSMTPPort, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		log.Fatal("failed to parse string to uint, err:", err.Error())
	}

	passwordResetTokenExpDuration, err := time.ParseDuration(os.Getenv("PASSWORD_RESET_TOKEN_EXPIRATION_DURATION"))
	if err != nil {
		log.Fatal("failed to parse string to time.Duration, err:", err.Error())
	}

	passwordResetTokenLength, err := strconv.Atoi(os.Getenv("PASSWORD_RESET_TOKEN_LENGTH"))
	if err != nil {
		log.Fatal("failed to parse string to uint, err:", err.Error())
	}

	return &appConfig{
		DB_DIALECT:                        os.Getenv("DB_DIALECT"),
		DB_HOST:                           os.Getenv("DB_HOST"),
		DB_PORT:                           uint(dbPort),
		DB_NAME:                           os.Getenv("DB_NAME"),
		DB_USER:                           os.Getenv("DB_USER"),
		DB_PASSWORD:                       os.Getenv("DB_PASSWORD"),
		APP_PORT:                          uint(appPort),
		ACCESS_TOKEN_SECRET:               os.Getenv("ACCESS_TOKEN_SECRET"),
		REFRESH_TOKEN_SECRET:              os.Getenv("REFRESH_TOKEN_SECRET"),
		ACCESS_TOKEN_EXPIRATION_DURATION:  accessTokenExpDuration,
		REFRESH_TOKEN_EXPIRATION_DURATION: refreshTokenExpDuration,
		ACCOUNT_ACTIVATION_TOKEN_EXPIRATION_DURATION: accountActivationTokenExpDuration,
		ACCOUNT_ACTIVATION_URL:                       os.Getenv("ACCOUNT_ACTIVATION_URL"),
		EMAIL_SMTP_SERVER:                            os.Getenv("EMAIL_SMTP_SERVER"),
		EMAIL_SMTP_PORT:                              uint(emailSMTPPort),
		EMAIL_SENDER_IDENTITY:                        os.Getenv("EMAIL_SENDER_IDENTITY"),
		EMAIL_SENDER_AND_SMTP_USER:                   os.Getenv("EMAIL_SENDER_AND_SMTP_USER"),
		EMAIL_SMTP_USER_PASSWORD:                     os.Getenv("EMAIL_SMTP_USER_PASSWORD"),
		PASSWORD_RESET_TOKEN_EXPIRATION_DURATION:     passwordResetTokenExpDuration,
		PASSWORD_RESET_TOKEN_LENGTH:                  uint(passwordResetTokenLength),
	}
}

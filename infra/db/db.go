package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
)

func InitDB() *sql.DB {
	appConfig := config.GetAppConfig()

	db, err := sql.Open(appConfig.DB_DIALECT, fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", appConfig.DB_HOST, appConfig.DB_PORT, appConfig.DB_NAME, appConfig.DB_USER, appConfig.DB_PASSWORD))

	if err != nil {
		log.Fatalf("failed to parse postgres dsn. (%s)", err.Error())
	}

	if errPing := db.Ping(); errPing != nil {
		log.Fatalf("failed to connect to db. (%s)", errPing.Error())
	}

	return db
}

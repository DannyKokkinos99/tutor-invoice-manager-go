package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	projectType := os.Getenv("PROJECT_TYPE")
	var dsn string
	if projectType == "development" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("DEV_POSTGRES_HOST"),
			os.Getenv("DEV_POSTGRES_USER"),
			os.Getenv("DEV_POSTGRES_PASSWORD"),
			os.Getenv("DEV_POSTGRES_DB"),
			os.Getenv("DEV_POSTGRES_PORT"),
		)
	} else if projectType == "production" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("PROD_POSTGRES_HOST"),
			os.Getenv("PROD_POSTGRES_USER"),
			os.Getenv("PROD_POSTGRES_PASSWORD"),
			os.Getenv("PROD_POSTGRES_DB"),
			os.Getenv("PROD_POSTGRES_PORT"),
		)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        dsn,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	return db, nil
}

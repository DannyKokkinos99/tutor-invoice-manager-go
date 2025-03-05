package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(projectType string) (*gorm.DB, error) {

	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}
	var dsn string
	var db *gorm.DB
	var err error
	if projectType == "development" {
		fmt.Println("IN DEV MODE")
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("DEV_POSTGRES_HOST"),
			os.Getenv("DEV_POSTGRES_USER"),
			os.Getenv("DEV_POSTGRES_PASSWORD"),
			os.Getenv("DEV_POSTGRES_DB"),
			os.Getenv("DEV_POSTGRES_PORT"),
		)
		db, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        dsn,
		}), &gorm.Config{})
		if err != nil {
			panic("failed to connect to database")
		}

	} else if projectType == "production" {
		fmt.Println("IN PRODUCTION")
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("PROD_POSTGRES_HOST"),
			os.Getenv("PROD_POSTGRES_USER"),
			os.Getenv("PROD_POSTGRES_PASSWORD"),
			os.Getenv("PROD_POSTGRES_DB"),
			os.Getenv("PROD_POSTGRES_PORT"),
		)
		db, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        dsn,
		}), &gorm.Config{
			PrepareStmt: false,
			Logger:      logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic("failed to connect to database")
		}

	} else {
		panic("Project type not defined:" + projectType)
	}

	return db, nil
}

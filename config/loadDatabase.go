package config

import (
	"fmt"
	"log"
	"os"

	"github.com/rafli-lutfi/go-commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type dbCreds struct {
	username     string
	password     string
	databaseName string
	port         string
	sslMode      string
	timeZone     string
}

func LoadDatabase() {
	var dbCreds = &dbCreds{
		username:     os.Getenv("DATABASE_USERNAME"),
		password:     os.Getenv("DATABASE_PASSWORD"),
		databaseName: os.Getenv("DATABASE_NAME"),
		port:         os.Getenv("DATABASE_PORT"),
		sslMode:      os.Getenv("DATABASE_SSLMODE"),
		timeZone:     os.Getenv("DATABASE_TIMEZONE"),
	}

	var err error

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbCreds.username,
		dbCreds.password,
		dbCreds.databaseName,
		dbCreds.port,
		dbCreds.sslMode,
		dbCreds.timeZone,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(
		models.Category{},
		models.Discount{},
		models.Product{},
	)

	fmt.Println("success connect to database ...")
}

func GetDBConnection() *gorm.DB {
	return db
}

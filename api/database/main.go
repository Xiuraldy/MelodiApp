package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func CreateDbConnection() {
	if testing.Testing() {
		var err error
		connStr := "host=localhost port=5432 user=postgres password=postgresql dbname=db_test sslmode=disable"
		DBConn, err = gorm.Open(postgres.Open(connStr), &gorm.Config{TranslateError: true})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL"))
		DBConn, err = gorm.Open(postgres.Open(connStr), &gorm.Config{TranslateError: true})
		if err != nil {
			log.Fatal(err)
		}
	}
}

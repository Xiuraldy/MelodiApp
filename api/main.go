package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"trucode.app/api/auth"
	"trucode.app/api/database"
	"trucode.app/api/models"
	"trucode.app/api/shared"
	"trucode.app/api/users"

	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(shared.Cors())
	users.AddUserRoutes(router)
	auth.AddAuthRoutes(router)

	return router
}

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Unable to load env vars")
		}
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	_, err := gorm.Open(postgres.Open(connStr))
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}

	router := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

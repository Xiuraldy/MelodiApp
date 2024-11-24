package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"trucode.app/api/auth"
	"trucode.app/api/census"
	"trucode.app/api/database"
	"trucode.app/api/models"
	"trucode.app/api/shared"
	"trucode.app/api/userConfig"
	"trucode.app/api/users"
)

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	requiredVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "PORT"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Error: missing required environment variable %s", v)
		}
	}
}

func setupDatabase() *gorm.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Person{}, &models.UserConfig{})
	if err != nil {
		log.Fatalf("Error during database migration: %v", err)
	}

	return db
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(shared.Cors())

	users.AddUserRoutes(router)
	census.AddCensusRoutes(router)
	auth.AddAuthRoutes(router)
	userConfig.AddUserConfigRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Success": true})
	})

	return router
}

func main() {
	loadEnvVars()

	database.DBConn = setupDatabase()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down gracefully...")
		os.Exit(0)
	}()

	router := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

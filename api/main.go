package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"trucode.app/api/auth"
	"trucode.app/api/database"
	"trucode.app/api/entries"
	"trucode.app/api/models"
	"trucode.app/api/shared"
	"trucode.app/api/users"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(shared.Cors())
	users.AddUserRoutes(router)
	entries.AddEntryRoutes(router)
	auth.AddAuthRoutes(router)

	return router
}

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	getData()
	loadEnvVars()
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{}, &models.Entry{})

	router := setupRouter()

	router.Run(":3333")
}

func getData() {
	file, err := os.Open("data/source.data")
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(file)

	csvReaderRead, _ := csvReader.ReadAll()
	for _, line := range csvReaderRead {
		fmt.Println(line)
	}
}

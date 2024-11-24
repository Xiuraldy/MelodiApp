package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"trucode.app/api/auth"
	"trucode.app/api/census"
	"trucode.app/api/database"
	"trucode.app/api/models"
	"trucode.app/api/shared"
	"trucode.app/api/userConfig"
	"trucode.app/api/users"

	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(shared.Cors())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Success": true})
	})

	users.AddUserRoutes(router)
	census.AddCensusRoutes(router)
	auth.AddAuthRoutes(router)
	userConfig.AddUserConfigRoutes(router)

	return router
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	database.CreateDbConnection()

	router := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down gracefully...")
		os.Exit(0)
	}()

	log.Printf("Server running on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func insertPerson(person models.Person, wg *sync.WaitGroup) {
	defer wg.Done()
	database.DBConn.Create(&person)
}

func getData() {

	file, err := os.Open("dataCensus/source.data")
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(file)

	csvReaderRead, _ := csvReader.ReadAll()

	var wg sync.WaitGroup

	for _, line := range csvReaderRead[1:] {
		wg.Add(1)

		age, _ := strconv.Atoi(strings.TrimSpace(line[0]))
		fnlwgt, _ := strconv.Atoi(strings.TrimSpace(line[2]))
		educationNum, _ := strconv.Atoi(strings.TrimSpace(line[4]))
		capitalGain, _ := strconv.Atoi(strings.TrimSpace(line[10]))
		capitalLoss, _ := strconv.Atoi(strings.TrimSpace(line[11]))
		hoursPerWeek, _ := strconv.Atoi(strings.TrimSpace(line[12]))

		person := models.Person{
			Age:           age,
			Workclass:     line[1],
			Fnlwgt:        fnlwgt,
			Education:     line[3],
			EducationNum:  educationNum,
			MaritalStatus: line[5],
			Occupation:    line[6],
			Relationship:  line[7],
			Race:          line[8],
			Sex:           line[9],
			CapitalGain:   capitalGain,
			CapitalLoss:   capitalLoss,
			HoursPerWeek:  hoursPerWeek,
			NativeCountry: line[13],
			Income:        line[14],
		}

		go insertPerson(person, &wg)

	}
	wg.Wait()
}

func clearData() {
	if err := database.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Person{}).Error; err != nil {
		log.Printf("Error al eliminar datos de Person: %v", err)
	}
	log.Println("Datos eliminados de la tabla Person.")
}

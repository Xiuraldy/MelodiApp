package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
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
	users.AddUserRoutes(router)
	census.AddCensusRoutes(router)
	auth.AddAuthRoutes(router)
	userConfig.AddUserConfigRoutes(router)

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
	database.DBConn.AutoMigrate(&models.User{}, &models.Person{}, &models.UserConfig{})

	DBConn, err := gorm.Open(postgres.Open(connStr))
	if err != nil {
		log.Fatal("Unable to connect to DB")
	}

	router := setupRouter()

	router.GET("/", func(c *gin.Context) {
		tx := DBConn.Exec("SELECT 1")
		if tx.Error != nil {
			log.Printf("Error al ejecutar consulta: %v", tx.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"Success": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Success": true})
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		clearData()
		os.Exit(0)
	}()

	getData()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	log.Printf("Server running on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
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

	csvReaderRead, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	if len(csvReaderRead) < 3 {
		log.Fatal("Error: csvReaderRead tiene menos de 3 elementos")
	}

	for _, line := range csvReaderRead[1:3] {
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

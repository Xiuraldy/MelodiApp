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

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(shared.Cors())
	users.AddUserRoutes(router)
	census.AddCensusRoutes(router)
	auth.AddAuthRoutes(router)
	userConfig.AddUserConfigRoutes(router)
	return router
}

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using system environment variables.")
	}

	requiredVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "PORT"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Error: missing required environment variable %s", v)
		}
	}
}

func main() {
	// Cargar las variables de entorno
	loadEnvVars()

	// Configurar la conexión a la base de datos
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	DBConn, err := gorm.Open(postgres.Open(connStr))
	if err != nil {
		log.Fatal("Unable to connect to DB: ", err)
	}
	database.DBConn = DBConn // Asigna la conexión globalmente

	// Migraciones de la base de datos
	if err := database.DBConn.AutoMigrate(&models.User{}, &models.Person{}, &models.UserConfig{}); err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	// Procesar los datos iniciales
	getData()

	// Manejo de señales para limpiar datos al salir
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		clearData()
		os.Exit(0)
	}()

	// Configuración del servidor
	router := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // Valor predeterminado
	}
	fmt.Printf("Server running on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func insertPerson(person models.Person, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := database.DBConn.Create(&person).Error; err != nil {
		log.Printf("Error inserting person: %v", err)
	}
}

func getData() {
	file, err := os.Open("dataCensus/source.data")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
		log.Printf("Error clearing data: %v", err)
	}
	log.Println("Data cleared from table Person.")
}

package census

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"trucode.app/api/auth"
	"trucode.app/api/database"
	"trucode.app/api/models"
	"trucode.app/api/shared"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(shared.Cors())
	AddCensusRoutes(router)
	auth.AddAuthRoutes(router)

	return router
}

func GetCensus(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	_, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var queryParams = map[string]string{
		"age":            c.Query("age"),
		"workclass":      c.Query("workclass"),
		"fnlwgt":         c.Query("fnlwgt"),
		"education":      c.Query("education"),
		"education_num":  c.Query("educationNum"),
		"marital_status": c.Query("maritalStatus"),
		"occupation":     c.Query("occupation"),
		"relationship":   c.Query("relationship"),
		"race":           c.Query("race"),
		"sex":            c.Query("sex"),
		"capital_gain":   c.Query("capitalGain"),
		"capital_loss":   c.Query("capitalLoss"),
		"hours_per_week": c.Query("hoursPerWeek"),
		"native_country": c.Query("nativeCountry"),
		"income":         c.Query("income"),
	}

	var people []models.Person
	query := database.DBConn

	for key, value := range queryParams {
		if value != "" {
			if strings.Contains(value, "to") {
				ranges := strings.Split(value, "to")
				lower := strings.TrimSpace(ranges[0])
				upper := strings.TrimSpace(ranges[1])
				query = query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", key), lower, upper)
			} else {
				query = query.Where(fmt.Sprintf("%s LIKE ?", key), "%"+value+"%")
			}
		}
	}

	var sortBy = c.Query("sortBy")
	var sortOrder = c.Query("sortOrder")

	if sortBy != "" && sortOrder != "" {
		if sortOrder == "asc" || sortOrder == "desc" {
			query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid sortOrder"})
			return
		}
	}

	tx := query.Find(&people)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not retrieve people: %v", tx.Error)})
		return
	}

	if len(people) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No entries found"})
		return
	}

	var paginator, errp = strconv.Atoi(c.Query("paginator"))
	if errp != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errp.Error()})
		return
	}

	var finishPag = paginator
	if finishPag > len(people) {
		finishPag = len(people)
	}
	var initialPag = paginator - 10

	var totalRecords int64
	tx = query.Model(&models.Person{}).Count(&totalRecords)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not count records: %v", tx.Error)})
		return
	}

	var census = people[initialPag:finishPag]

	roundedTotalRecords := int64(math.Ceil(float64(totalRecords)))

	c.JSON(http.StatusOK, gin.H{
		"total":   roundedTotalRecords,
		"current": census,
	})
}

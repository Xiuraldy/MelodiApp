package userConfig

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"trucode.app/api/database"
	"trucode.app/api/models"
	"trucode.app/api/shared"
)

func getUserConfig(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(*shared.Payload)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userId, ok := claims.MapClaims["user_id"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user_id not found in token"})
		return
	}

	var userConfig models.UserConfig
	database.DBConn.Where("user_id=?", userId).First(&userConfig)

	fmt.Println(userConfig)

	c.JSON(http.StatusOK, userConfig)
}

func updateUserConfig(c *gin.Context) {
	var userConfigActual *models.UserConfig
	var userInput models.UserConfig

	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	claims, ok := token.Claims.(*shared.Payload)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userId, ok := claims.MapClaims["user_id"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user_id not found in token"})
		return
	}

	tx := database.DBConn.Where("user_id=?", int(userId)).First(&userConfigActual)
	fmt.Println("tx", int(userId))
	fmt.Println("tx", tx.Error)

	userConfigActual.UserID = int(userId)
	userConfigActual.Age = userInput.Age
	userConfigActual.Workclass = userInput.Workclass
	userConfigActual.Fnlwgt = userInput.Fnlwgt
	userConfigActual.Education = userInput.Education
	userConfigActual.EducationNum = userInput.EducationNum
	userConfigActual.MaritalStatus = userInput.MaritalStatus
	userConfigActual.Occupation = userInput.Occupation
	userConfigActual.Relationship = userInput.Relationship
	userConfigActual.Race = userInput.Race
	userConfigActual.Sex = userInput.Sex
	userConfigActual.CapitalGain = userInput.CapitalGain
	userConfigActual.CapitalLoss = userInput.CapitalLoss
	userConfigActual.HoursPerWeek = userInput.HoursPerWeek
	userConfigActual.NativeCountry = userInput.NativeCountry
	userConfigActual.Income = userInput.Income
	userConfigActual.SortBy = userInput.SortBy
	userConfigActual.SortOrder = userInput.SortOrder
	userConfigActual.Paginator = userInput.Paginator

	fmt.Println("userInput", userInput)
	fmt.Println("userConfig1", userConfigActual)

	if tx := database.DBConn.Save(&userConfigActual); tx.Error != nil {
		fmt.Println("tx.Error", tx.Error)
		if tx := database.DBConn.Create(&userConfigActual); tx.Error != nil {
			fmt.Println("Entra", tx.Error)
			if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, gin.H{"error": "userConfig already exists"})
				return
			}
		}
	}

	fmt.Println("userConfig", userConfigActual)

	c.JSON(http.StatusOK, gin.H{"userConfigUpdate": userConfigActual})
}

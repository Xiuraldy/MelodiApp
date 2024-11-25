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
	if tx := database.DBConn.Where("user_id=?", userId).First(&userConfig); tx.Error != nil {
		if tx.Error.Error() == gorm.ErrRecordNotFound.Error() {
			fmt.Println("-----> No hay registros")
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		fmt.Println(tx)
		c.JSON(http.StatusConflict, gin.H{"error": tx.Error})
		return
	}

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

	userInput.UserID = int(userId)

	tx := database.DBConn.Where("user_id=?", int(userId)).First(&userConfigActual)
	if tx.Error != nil {
		if tx := database.DBConn.Create(&userInput); tx.Error != nil {
			fmt.Println("entra")
			if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, gin.H{"error": "userConfig already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"userConfig": userInput})
		return
	}

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

	fmt.Println("userConfigActual", userConfigActual)

	if tx := database.DBConn.Model(&userConfigActual).
		Where("user_id = ?", int(userId)).
		Select("*").
		Updates(userConfigActual); tx.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": tx.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userConfigUpdate": userConfigActual})
}

package users

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"trucode.app/api/auth"
	"trucode.app/api/census"
	"trucode.app/api/database"
	"trucode.app/api/models"
	"trucode.app/api/shared"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(shared.Cors())
	AddUserRoutes(router)
	census.AddCensusRoutes(router)
	auth.AddAuthRoutes(router)

	return router
}

func GetAllUsers(c *gin.Context) {
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

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var user models.User

	tx := database.DBConn.Where("id=?", userData.Uid).Find(&user)
	if tx.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var users []models.User
	database.DBConn.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {
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

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	if userData.ExpiryTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var user models.User

	tx := database.DBConn.Where("id=?", userData.Uid).Find(&user)
	if tx.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	tx = database.DBConn.Where("id=?", c.Param("id")).Find(&user)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetMe(c *gin.Context) {
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

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	if userData.ExpiryTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var foundUser models.User
	tx := database.DBConn.Where("id=?", userData.Uid).Find(&foundUser)
	if tx.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

func CreateUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	database.DBConn.Create(&user)
	c.JSON(http.StatusCreated, user)
}

func EditUser(c *gin.Context) {
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

	claims, _ := token.Claims.(*shared.Payload)
	_, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")

	var user models.User
	tx := database.DBConn.First(&user, id)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updatedData models.User
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user.Username = updatedData.Username
	user.Email = updatedData.Email
	if updatedData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	database.DBConn.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	database.DBConn.Delete(&models.User{}, id)
	c.JSON(http.StatusNoContent, gin.H{"id": id})
}

package entries

import (
	"fmt"
	"net/http"
	"os"

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
	AddEntryRoutes(router)
	auth.AddAuthRoutes(router)

	return router
}

func GetAllEntries(c *gin.Context) {
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
	entryData, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var user models.Entry

	tx := database.DBConn.Where("id=?", entryData.Uid).Find(&user)
	if tx.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var entries []models.Entry
	database.DBConn.Find(&entries)
	c.JSON(http.StatusOK, entries)
}

func GetMyEntries(c *gin.Context) {
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
	entryData, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var foundEntries []models.Entry

	tx := database.DBConn.Where("user_id = ?", entryData.Uid).Find(&foundEntries)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not retrieve entries: %v", tx.Error)})
		return
	}

	if len(foundEntries) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No entries found"})
		return
	}

	c.JSON(http.StatusOK, foundEntries)
}

func CreateEntry(c *gin.Context) {
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
	entryData, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	var entry models.Entry
	c.BindJSON(&entry)

	entry.UserID = entryData.Uid

	database.DBConn.Create(&entry)
	c.JSON(http.StatusCreated, entry)
}

func DeleteEntry(c *gin.Context) {
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
	entryData, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")

	var entry models.Entry
	tx := database.DBConn.First(&entry, id)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}
	if entry.UserID != entryData.Uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this entry"})
		return
	}

	database.DBConn.Delete(&entry)
	c.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully", "id": id})
}

func EditEntry(c *gin.Context) {
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
	entryData, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")

	var entry models.Entry
	tx := database.DBConn.First(&entry, id)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}
	if entry.UserID != entryData.Uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this entry"})
		return
	}

	var updatedData models.Entry
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	entry.Title = updatedData.Title
	entry.Content = updatedData.Content

	database.DBConn.Save(&entry)

	c.JSON(http.StatusOK, gin.H{"message": "Entry updated successfully", "entry": entry})
}

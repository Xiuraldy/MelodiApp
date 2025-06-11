package shared

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"trucode.app/api/database"
	"trucode.app/api/models"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Private-Network", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthenticateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := GetTokenFromRequest(c)
		token, err := jwt.ParseWithClaims(tokenStr, &Payload{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid token")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, _ := token.Claims.(*Payload)

		userData := Sessions[claims.Session]
		if userData.ExpiryTime.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired session"})
			return
		}

		var user models.User
		tx := database.DBConn.Where("id=?", userData.Uid).Find(&user)
		if tx.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": tx.Error.Error()})
			return
		}

		c.Set("authorizedUser", user)

		c.Next()
	}
}

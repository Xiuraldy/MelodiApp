package shared

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	jwt.MapClaims        // ExpiryAt, IssueAt
	Session       string `json:"session"`
}

func GetTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	tokenSlice := strings.Split(bearerToken, " ")
	if len(tokenSlice) == 2 {
		return tokenSlice[1]
	}
	return ""
}

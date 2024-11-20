package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"trucode.app/api/database"
	"trucode.app/api/models"
	"trucode.app/api/shared"
)

func Register(c *gin.Context) {
	var userInput models.UserInput

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Username: userInput.Username,
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	if userInput.Username == "" || userInput.Email == "" || userInput.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incomplete fields"})
		return
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(userInput.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	var existingUser models.User
	if err := database.DBConn.Where("email = ?", userInput.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	if tx := database.DBConn.Create(&user); tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
	}

	sessionToken := uuid.NewV5(uuid.UUID{}, "session").String()

	session := shared.Session{
		Uid:        user.ID,
		ExpiryTime: time.Now().Add(10 * time.Minute),
	}

	shared.Sessions[sessionToken] = session

	claims := shared.Payload{
		MapClaims: jwt.MapClaims{
			"iat":     jwt.NewNumericDate(time.Now()),
			"eat":     jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			"user_id": user.ID,
		},
		Session: sessionToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signinKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString, err := token.SignedString(signinKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Login(c *gin.Context) {
	var input models.UserInput
	var user models.User

	c.BindJSON(&input)

	if input.Email == "" || input.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incomplete fields"})
		return
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	database.DBConn.Where("email=?", input.Email).Find(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	sessionToken := uuid.NewV5(uuid.UUID{}, "session").String()

	session := shared.Session{
		Uid:        user.ID,
		ExpiryTime: time.Now().Add(10 * time.Minute),
	}

	shared.Sessions[sessionToken] = session

	claims := shared.Payload{
		MapClaims: jwt.MapClaims{
			"iat":     jwt.NewNumericDate(time.Now()),
			"eat":     jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			"user_id": user.ID,
		},
		Session: sessionToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signinKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString, err := token.SignedString(signinKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
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

	delete(shared.Sessions, claims.Session)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

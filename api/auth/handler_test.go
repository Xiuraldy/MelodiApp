package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"trucode.app/api/database"
	"trucode.app/api/models"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.DELETE("/logout", Logout)
	return r
}

func getAuthToken(t *testing.T, router http.Handler) string {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	user := models.User{
		Username: "xiu",
		Email:    "x@gmail.com",
		Password: "12345",
	}
	jsonValue, _ := json.Marshal(user)

	wCreate := httptest.NewRecorder()
	reqCreate := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(wCreate, reqCreate)

	user = models.User{
		Email:    "x@gmail.com",
		Password: "12345",
	}
	jsonValue, _ = json.Marshal(user)

	wAuth := httptest.NewRecorder()
	reqAuth := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(wAuth, reqAuth)

	responseParts := strings.Split(wAuth.Body.String(), "}{")
	if len(responseParts) > 1 {
		responseParts[0] = responseParts[0] + "}"
		responseParts[1] = "{" + responseParts[1]
	}

	var loginResponse struct {
		Token string `json:"token"`
	}

	if err := json.Unmarshal([]byte(responseParts[0]), &loginResponse); err != nil {
		t.Fatalf("Error parsing JSON response: %v", err)
	}

	return loginResponse.Token
}

func TestRegisterSuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	defer database.DBConn.Migrator().DropTable(&models.User{}) // Cleanup after test

	router := setupRouter()

	userInput := models.UserInput{
		Username: "xiu",
		Email:    "new@gmail.com",
		Password: "12345",
	}

	jsonValue, _ := json.Marshal(userInput)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "xiu")
	clearDatabase()
}

func TestRegisterEmailConflict(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	defer database.DBConn.Migrator().DropTable(&models.User{}) // Cleanup after test

	router := setupRouter()

	user := models.User{
		Username: "andres",
		Email:    "and@gmail.com",
		Password: "12345",
	}
	database.DBConn.Create(&user)

	userInput := models.UserInput{
		Username: "wences",
		Email:    "and@gmail.com",
		Password: "12345",
	}

	jsonValue, _ := json.Marshal(userInput)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Email already exists")
	clearDatabase()
}

func TestLoginSuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	defer database.DBConn.Migrator().DropTable(&models.User{}) // Cleanup after test

	router := setupRouter()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	user := models.User{
		Username: "xiu",
		Email:    "x@gmail.com",
		Password: string(hashedPassword),
	}
	database.DBConn.Create(&user)

	userInput := models.UserInput{
		Email:    "x@gmail.com",
		Password: string(hashedPassword),
	}

	jsonValue, _ := json.Marshal(userInput)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
	clearDatabase()
}

func TestLoginInvalidCredentials(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	defer database.DBConn.Migrator().DropTable(&models.User{}) // Cleanup after test

	router := setupRouter()

	userInput := models.UserInput{
		Email:    "x@gmail.com",
		Password: "banana",
	}

	jsonValue, _ := json.Marshal(userInput)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
	clearDatabase()
}

func TestLogoutSuccess(t *testing.T) {
	router := setupRouter()
	token := getAuthToken(t, router)
	fmt.Println("token", token)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/logout", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Successfully logged out")
}

func TestLogoutUnauthorized(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/logout", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func clearDatabase() {
	database.DBConn.Exec("DELETE FROM users")
	database.DBConn.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

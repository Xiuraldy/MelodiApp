package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"trucode.app/api/database"
	"trucode.app/api/models"
)

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
	reqCreate := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(wCreate, reqCreate)

	user = models.User{
		Email:    "x@gmail.com",
		Password: "12345",
	}
	jsonValue, _ = json.Marshal(user)

	wAuth := httptest.NewRecorder()
	reqAuth := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonValue))
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

func TestGetAllUsersUnauthorized(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateUserSuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})
	router := setupRouter()

	user := models.User{
		Username: "Wences",
		Email:    "wences@gmail.com",
		Password: "12345",
	}
	jsonValue, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "wences@gmail.com")
	clearDatabase()
}

func TestGetUserByIdNotFound(t *testing.T) {
	router := setupRouter()

	token := getAuthToken(t, router)

	wGetUserById := httptest.NewRecorder()
	reqGetUserById := httptest.NewRequest("GET", "/users/222", nil)
	reqGetUserById.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(wGetUserById, reqGetUserById)

	assert.Equal(t, http.StatusNotFound, wGetUserById.Code)
	assert.Contains(t, wGetUserById.Body.String(), "User not found")
	clearDatabase()
}

func TestEditUserSuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})
	router := setupRouter()

	user := models.User{Username: "xiu", Email: "x@gmail.com", Password: "12345"}
	database.DBConn.Create(&user)

	updatedUser := models.User{
		Username: "xiuraldy",
		Email:    "xiu@gmail.com",
		Password: "123",
	}
	jsonValue, _ := json.Marshal(updatedUser)

	token := getAuthToken(t, router)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "xiu@gmail.com")
	clearDatabase()
}

func TestDeleteUserSuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})
	router := setupRouter()

	user := models.User{Username: "xiu", Email: "x@gmail.com", Password: "12345"}
	database.DBConn.Create(&user)

	token := getAuthToken(t, router)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%d", user.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	clearDatabase()
}

func clearDatabase() {
	database.DBConn.Exec("DELETE FROM users")
	database.DBConn.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

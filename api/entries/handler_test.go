package entries

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

func TestGetAllEntriesUnauthorized(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entries", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateEntrySuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.Entry{})

	router := setupRouter()

	entry := models.Entry{
		Title:   "Title Entry",
		Content: "This is a content.",
	}
	jsonValue, _ := json.Marshal(entry)

	token := getAuthToken(t, router)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/entries", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Title Entry")
	clearDatabaseEntries()
}

func TestGetMyEntriesSuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.Entry{})

	router := setupRouter()
	token := getAuthToken(t, router)

	entry := models.Entry{
		Title:   "Title my Entry",
		Content: "This is a content.",
		UserID:  1,
	}
	database.DBConn.Create(&entry)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/entries/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Title my Entry")
	clearDatabaseEntries()
}

func TestDeleteEntrySuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.Entry{})

	router := setupRouter()
	token := getAuthToken(t, router)

	entry := models.Entry{
		Title:   "Title entry",
		Content: "This is a content.",
		UserID:  1,
	}
	database.DBConn.Create(&entry)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/entries/%d", entry.ID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Entry deleted successfully")
	clearDatabaseEntries()
}

func TestEditEntrySuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.Entry{})

	router := setupRouter()
	token := getAuthToken(t, router)

	entry := models.Entry{
		Title:   "Entry",
		Content: "content",
		UserID:  1,
	}
	database.DBConn.Create(&entry)

	updatedEntry := models.Entry{
		Title:   "New Entry",
		Content: "New content",
	}
	jsonValue, _ := json.Marshal(updatedEntry)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", fmt.Sprintf("/entries/%d", entry.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "New Entry")
	clearDatabaseEntries()
}

func clearDatabaseEntries() {
	database.DBConn.Exec("DELETE FROM entries")
	database.DBConn.Exec("ALTER SEQUENCE entries_id_seq RESTART WITH 1")
}

package census

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

func TestCreatePersonSuccess(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.Person{})

	router := setupRouter()

	person := models.Person{
		Age:           25,
		Workclass:     "Private",
		Fnlwgt:        226802,
		Education:     "Bachelors",
		EducationNum:  13,
		MaritalStatus: "Never-married",
		Occupation:    "Tech-support",
		Relationship:  "Not-in-family",
		Race:          "White",
		Sex:           "Male",
		CapitalGain:   0,
		CapitalLoss:   0,
		HoursPerWeek:  40,
		NativeCountry: "United-States",
		Income:        ">50K",
	}
	jsonValue, _ := json.Marshal(person)

	token := getAuthToken(t, router)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/census", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Age person")
	clearDatabaseEntries()
}

// func TestGetMyEntriesSuccess(t *testing.T) {
// 	database.CreateDbConnection()
// 	database.DBConn.AutoMigrate(&models.Person{})

// 	router := setupRouter()
// 	token := getAuthToken(t, router)

// 	entry := models.Entry{
// 		Title:   "Title my Entry",
// 		Content: "This is a content.",
// 		UserID:  1,
// 	}
// 	database.DBConn.Create(&entry)

// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest("GET", "/entries/me", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "Title my Entry")
// 	clearDatabaseEntries()
// }

func clearDatabaseEntries() {
	database.DBConn.Exec("DELETE FROM entries")
	database.DBConn.Exec("ALTER SEQUENCE entries_id_seq RESTART WITH 1")
}

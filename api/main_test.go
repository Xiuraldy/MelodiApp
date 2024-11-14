package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"trucode.app/api/database"
	"trucode.app/api/models"
)

func TestRegisterBadRequest(t *testing.T) {
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	router := setupRouter()     // Configurar las rutas
	w := httptest.NewRecorder() // Recorre las respuestas

	req := httptest.NewRequest("POST", "/auth/register", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestRegisterSuccess(t *testing.T) {
	loadEnvVars()
	database.CreateDbConnection()
	database.DBConn.AutoMigrate(&models.User{})

	router := setupRouter()     // Configurar las rutas
	w := httptest.NewRecorder() // Recorre las respuestas
	input := models.UserInput{
		Username: "andres",
		Email:    "and@gmail.com",
		Password: "12345",
	}
	jsonVal, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonVal))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "{\"email\":\"and@gmail.com\",\"id\":3,\"username\":\"andres\"}", w.Body.String())
}

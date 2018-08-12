package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type currencyRequest struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPing(t *testing.T) {
	// Build our expected body
	body := gin.H{
		"ping": "pong",
	}
	// Grab our router
	router := SetupRouter()
	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/ping")
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["ping"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["ping"], value)
}

func TestListCurrency(t *testing.T) {
	router := SetupRouter()
	// Perform a GET request with that handler.
	w := performRequest(router, "GET", "/api/v1/currency/list")
	assert.Equal(t, 200, w.Code, "OK response is expected")
}

func TestAddCurrency(t *testing.T) {
	// Grab our router
	router := SetupRouter()
	currency := &currencyRequest{
		From: "SGD",
		To:   "IDR",
	}
	jsonCurrency, _ := json.Marshal(currency)
	request, _ := http.NewRequest("POST", "/api/v1/currency/", bytes.NewBuffer(jsonCurrency))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

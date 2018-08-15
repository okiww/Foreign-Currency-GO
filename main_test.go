package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

type response struct {
	data   []currencies
	status string
}

type currencies struct {
	id    int
	from  string
	to    string
	rates []rates
}

type rates struct {
	id   int
	date string
	rate string
}

func TestPing(t *testing.T) {

	r := gofight.New()

	r.GET("/ping").
		// turn on the debug mode.
		SetDebug(true).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {

			assert.Equal(t, "pong", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestApiListCurrency(t *testing.T) {

	r := gofight.New()

	r.GET("/api/v1/currency/list").
		// turn on the debug mode.
		SetDebug(true).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var resp response
			temp, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(temp, &resp)
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
func TestApiDeleteCurrency(t *testing.T) {

	r := gofight.New()

	r.DELETE("/api/v1/currency/delete").SetJSON(gofight.D{
		"currency_id": 1,
	}).Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		if r.Code == 400 {
			//if data not found
			assert.Equal(t, http.StatusBadRequest, r.Code)
		} else {
			//if success delete
			assert.Equal(t, http.StatusOK, r.Code)
		}
		assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
	})
}
func TestApiAddCurrency(t *testing.T) {

	r := gofight.New()

	r.POST("/api/v1/currency/add").SetJSON(gofight.D{
		"from": "JPY",
		"to":   "SGD",
	}).Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		if r.Code == 400 {
			//if bad request
			assert.Equal(t, http.StatusBadRequest, r.Code)
		} else {
			//if success add currency
			assert.Equal(t, http.StatusOK, r.Code)
		}
		assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
	})
}

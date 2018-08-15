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
		} else if r.Code == 200 {
			//if currency already in db
			assert.Equal(t, http.StatusOK, r.Code)
		} else {
			//if success add currency
			assert.Equal(t, http.StatusCreated, r.Code)
		}
		assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
	})
}

func TestApiAddRate(t *testing.T) {

	r := gofight.New()

	r.POST("/api/v1/rate/add").SetJSON(gofight.D{
		"currency_id": 7,
		"rate":        0.009,
		"date":        "2017-09-09",
	}).Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		if r.Code == 400 {
			//if bad request
			assert.Equal(t, http.StatusBadRequest, r.Code)
		} else {
			//if success add rate
			assert.Equal(t, http.StatusCreated, r.Code)
		}
		assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
	})
}

func TestApiListRateByDate(t *testing.T) {

	r := gofight.New()
	param := "2017-01-01"
	r.GET("/api/v1/rate/list?date="+param).
		// turn on the debug mode.
		SetDebug(true).
		Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestAPIGetMost7DataPointByCurrency(t *testing.T) {

	r := gofight.New()

	r.POST("/api/v1/rate/most").SetJSON(gofight.D{
		"currency_id": 7,
	}).Run(SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		if r.Code == 400 {
			//if bad request
			assert.Equal(t, http.StatusBadRequest, r.Code)
		} else {
			//if success add rate
			assert.Equal(t, http.StatusOK, r.Code)
		}
		assert.Equal(t, "application/json; charset=utf-8", r.HeaderMap.Get("Content-Type"))
	})
}

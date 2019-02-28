package main

import (
	"log"
	"myblog/app"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", app.HomeHandler).Methods("GET")
	router.HandleFunc("/posts/{slug}", app.PostHandler).Methods("GET")
	router.PathPrefix("/").HandlerFunc(app.CatchAllHandler)
	return router
}

func TestHomeHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
	assert.Equal(t, "text/html; charset=utf-8", response.Result().Header["Content-Type"][0], "http content-type header response is expected")
}

func TestPostHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "/posts/a-non-existant-post", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "404 response is expected")
	log.Print(response.Result().Header)
	// assert.Equal(t, "text/html; charset=utf-8", response.Result().Header["Content-Type"][0], "http content-type header response is expected")

	request, _ = http.NewRequest("GET", "/posts/2019-02-26-website-in-a-binary", nil)
	response = httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "200 response is expected")
}

func TestAboutHandler(t *testing.T) {
}

func TestCatchAllHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "/foo", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "404 response is expected")
}

package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", homeEndpoint).Methods("GET")
	router.HandleFunc("/posts/{slug}", postEndpoint).Methods("GET")
	return router
}

func TestHomeEndpoint(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
	assert.Equal(t, "text/html; charset=utf-8", response.Result().Header["Content-Type"][0], "http content-type header response is expected")
}

func TestPostEndpoint(t *testing.T) {
	request, _ := http.NewRequest("GET", "/posts/a-non-existant-post", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "404 response is expected")
	log.Print(response.Result().Header)
	// assert.Equal(t, "text/html; charset=utf-8", response.Result().Header["Content-Type"][0], "http content-type header response is expected")

	request, _ = http.NewRequest("GET", "/posts/test_post_2", nil)
	response = httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "200 response is expected")
}

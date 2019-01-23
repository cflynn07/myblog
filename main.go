package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func homeEndpoint(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
}

func main() {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(path) + "/web/static/"
	router := mux.NewRouter()

	router.HandleFunc("/", homeEndpoint).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	log.Println("Listening...")
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

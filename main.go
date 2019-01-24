package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type pageProperties struct {
	Title string
}

func homeEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("homeEndpoint")
	pp := pageProperties{"Home"}
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/index.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.ExecuteTemplate(w, "layout", pp)
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

	log.Println("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

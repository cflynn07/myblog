package main

import (
	"fmt"
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

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/home.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.ExecuteTemplate(w, "layout", pp)
	if err != nil {
		log.Println(err)
	}
}

func postEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "post")
}

func aboutEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "about")
}

func contactEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "contact")
}

func main() {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(path) + "/web/static/"
	log.Println("dir: " + dir)

	router := mux.NewRouter()
	router.HandleFunc("/post/{slug}", postEndpoint).Methods("GET")
	router.HandleFunc("/about", aboutEndpoint).Methods("GET")
	router.HandleFunc("/contact", contactEndpoint).Methods("GET")
	router.HandleFunc("/", homeEndpoint).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	log.Println("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

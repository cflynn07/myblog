package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
}

func main() {
	router := httprouter.New()

	router.GET("/", handler)
	router.ServeFiles("/static/*filepath", http.Dir("/web/"))

	log.Println("Listening...")
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

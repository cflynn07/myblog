package main

import (
	"fmt"
	"log"
	"myblog/app"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
)

var path, _ = os.Executable()
var baseDir = filepath.Dir(path)
var staticDir = baseDir + "/web/static/"

var staticBox = packr.New("Static", "./web/static")

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})

	// Redirect to HTTPS in prod
	sslRedirect := (os.Getenv("SSL_REDIRECT") == "true")
	secureMiddleware := secure.New(secure.Options{
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		SSLRedirect:     sslRedirect,
	})

	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.Use(secureMiddleware.Handler)
	subRouter.HandleFunc("/", app.HomeHandler).Methods("GET")
	subRouter.HandleFunc("/about", app.AboutHandler).Methods("GET")
	subRouter.HandleFunc("/posts/{slug}", app.PostHandler).Methods("GET")
	subRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticBox)))

	router.PathPrefix("/").HandlerFunc(app.CatchAllHandler)

	log.Println("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

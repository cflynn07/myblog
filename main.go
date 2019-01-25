package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type pageProperties struct {
	Title       string
	PostContent template.HTML
}

var posts []string
var baseDir string
var staticDir string
var postsDir string

func homeEndpoint(w http.ResponseWriter, r *http.Request) {
	var input []byte
	input = []byte{1}
	log.Println("homeEndpoint")
	output := blackfriday.Run(input)
	log.Println(output)

	pp := pageProperties{
		Title:       "Home",
		PostContent: template.HTML(""),
	}

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
	vars := mux.Vars(r)
	log.Println(vars["slug"])
	log.Println(posts)
	data, err := ioutil.ReadFile(postsDir + vars["slug"] + ".md")
	if err != nil {
		log.Fatal(err)
	}
	output := blackfriday.Run(data)
	log.Println("output-----")
	log.Println(output)
	pp := pageProperties{
		Title:       "Home",
		PostContent: template.HTML(output),
	}
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/home.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.ExecuteTemplate(w, "layout", pp)
	if err != nil {
		log.Println(err)
	}
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
	baseDir = filepath.Dir(path)
	staticDir = baseDir + "/web/static/"
	postsDir = baseDir + "/web/posts/"

	// Get contents of web/posts, each file is a valid post route (/posts/{post})
	files, err := ioutil.ReadDir(postsDir)
	if err != nil {
		log.Fatal(err)
	}
	posts = make([]string, len(files))
	for i, file := range files {
		posts[i] = file.Name()
	}

	router := mux.NewRouter()
	router.HandleFunc("/post/{slug}", postEndpoint).Methods("GET")
	router.HandleFunc("/about", aboutEndpoint).Methods("GET")
	router.HandleFunc("/contact", contactEndpoint).Methods("GET")
	router.HandleFunc("/", homeEndpoint).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	log.Println("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

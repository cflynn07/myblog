package main

import (
	"fmt"
	"github.com/gorilla/mux"
	// "github.com/russross/blackfriday/v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Template variables for all pages
type globalPageVars struct {
	Title string
}

// Metadata for each blog post
type postMetaData struct {
	Title    string
	Subtitle string
	Keywords string
	Date     string // time.Time?
}

// All published blog posts
type blogPosts map[string]postMetaData

var gpv = globalPageVars{
	Title: "Casey Flynn",
}

// Blog uses this static map to display blog posts, the keys should match files
// in the template folder (sans the extension)
var bp = blogPosts{
	"test_post_2": postMetaData{
		Title:    "Test Post 2",
		Subtitle: "",
		Keywords: "",
		Date:     "",
	},
	"unicode_and_utf8": postMetaData{
		Title:    "Unicode and UTF8",
		Subtitle: "",
		Keywords: "",
		Date:     "",
	},
}

var posts []string
var baseDir string
var staticDir string
var postsDir string

func homeEndpoint(w http.ResponseWriter, r *http.Request) {
	type homePageVars struct {
		globalPageVars
		SubTitle  string
		BlogPosts blogPosts
	}

	hpv := homePageVars{
		globalPageVars: gpv,
		SubTitle:       "",
		BlogPosts:      bp,
	}

	t := template.Must(template.New("").Funcs(template.FuncMap{
		"trunc": func(s string) string {
			return s + "--test"
		},
	}).ParseFiles("templates/layout.html", "templates/home.html"))
	err := t.ExecuteTemplate(w, "layout", hpv)
	if err != nil {
		log.Print(err)
	}
}

/*
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
	pp := pageVars{
		Title: "Home",
		// PostContent: template.HTML(output),
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
*/
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
	postsDir = baseDir + "/posts/"

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
	// router.HandleFunc("/post/{slug}", postEndpoint).Methods("GET")
	router.HandleFunc("/about", aboutEndpoint).Methods("GET")
	router.HandleFunc("/contact", contactEndpoint).Methods("GET")
	router.HandleFunc("/", homeEndpoint).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	log.Println("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

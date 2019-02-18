package main

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Template variables for all pages
type globalPageVars struct {
	Title string
}

// Metadata for each blog post
type postData struct {
	Title          string
	Subtitle       string
	Keywords       string
	Date           string // time.Time?
	Content        string
	ContentPreview string
}

// All published blog posts
type blogPosts map[string]*postData

var gpv = globalPageVars{
	Title: "Casey Flynn",
}

// Blog uses this static map to display blog posts, the keys should match files
// in the template folder (sans the extension)
var bp = blogPosts{
	"test_post_2": &postData{
		Title:    "Test Post 2",
		Subtitle: "",
		Keywords: "",
		Date:     "",
	},
	"unicode_and_utf8": &postData{
		Title:    "Unicode and UTF8",
		Subtitle: "",
		Keywords: "",
		Date:     "",
	},
}

var box = packr.New("Posts", "./posts")

var path, _ = os.Executable()
var baseDir = filepath.Dir(path)
var staticDir = baseDir + "/web/static/"

func truncHelper(s string) string {
	words := strings.Fields(s)
	maxPreviewLength := int(math.Min(40, float64(len(words))))
	words = words[0:maxPreviewLength]
	return strings.Join(words, " ")
}

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

	for key, value := range bp {
		// data, err := ioutil.ReadFile(postsDir + key + ".md")
		data, err := box.Find(key + ".md")
		if err != nil {
			log.Fatal("error opening post", err)
		}
		value.Content = string(blackfriday.Run(data))
		value.ContentPreview = truncHelper(string(data))
		log.Print("=================")
		log.Print(value.ContentPreview)
	}

	log.Print(bp["test_post_2"].ContentPreview)

	t := template.Must(template.New("").ParseFiles("templates/layout.html", "templates/home.html"))
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
	router := mux.NewRouter()
	// router.HandleFunc("/post/{slug}", postEndpoint).Methods("GET")
	router.HandleFunc("/about", aboutEndpoint).Methods("GET")
	router.HandleFunc("/contact", contactEndpoint).Methods("GET")
	router.HandleFunc("/", homeEndpoint).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	log.Println("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

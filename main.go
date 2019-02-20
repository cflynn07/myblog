package main

import (
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
	Host  string
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
	Host:  "https://cflynn.us",
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

var postsBox = packr.New("Posts", "./posts")
var staticBox = packr.New("Static", "./web/static")
var templateBox = packr.New("Templates", "./templates")

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
		Path      string
	}

	hpv := homePageVars{
		globalPageVars: gpv,
		SubTitle:       "",
		BlogPosts:      bp,
		Path:           "/",
	}

	for key, value := range bp {
		// data, err := ioutil.ReadFile(postsDir + key + ".md")
		data, err := postsBox.Find(key + ".md")
		if err != nil {
			log.Fatal("error opening post", err)
		}
		value.Content = string(blackfriday.Run(data))
		value.ContentPreview = truncHelper(string(data))
	}

	templateLayout, err := templateBox.FindString("layout.html")
	if err != nil {
		log.Fatal(err)
	}
	templateHome, err := templateBox.FindString("home.html")
	if err != nil {
		log.Fatal(err)
	}
	t := template.New("")
	t.Parse(templateLayout)
	t.Parse(templateHome)

	err = t.ExecuteTemplate(w, "layout", hpv)
	if err != nil {
		log.Print(err)
	}
}

func postEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	type postPageVars struct {
		globalPageVars
		SubTitle     string
		BlogPost     template.HTML
		BlogPostMeta *postData
		Path         string
	}

	ppv := postPageVars{
		globalPageVars: gpv,
		SubTitle:       "",
		Path:           "/posts/" + vars["slug"],
	}

	templateLayout, err := templateBox.FindString("layout.html")
	if err != nil {
		log.Fatal(err)
	}
	var templateContent string

	if post, ok := bp[vars["slug"]]; !ok {
		// Post slug doesn't exist
		w.WriteHeader(http.StatusNotFound)
		// use 404 template
		templateContent, err = templateBox.FindString("404.html")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		ppv.BlogPostMeta = post
		data, err := postsBox.Find(vars["slug"] + ".md")
		if err != nil {
			log.Fatal(err)
		}
		blogPost := blackfriday.Run(data)
		ppv.BlogPost = template.HTML(blogPost)
		// use post template
		templateContent, err = templateBox.FindString("post.html")
		if err != nil {
			log.Fatal(err)
		}
	}

	t := template.New("")
	t.Parse(templateLayout)
	t.Parse(templateContent)
	err = t.ExecuteTemplate(w, "layout", ppv)
	if err != nil {
		log.Println(err)
	}
}

func aboutEndpoint(w http.ResponseWriter, r *http.Request) {
	type aboutPageVars struct {
		globalPageVars
		Path string
	}

	apv := aboutPageVars{
		globalPageVars: gpv,
		Path:           "/about",
	}

	templateLayout, err := templateBox.FindString("layout.html")
	if err != nil {
		log.Fatal(err)
	}
	templateContent, err := templateBox.FindString("about.html")
	if err != nil {
		log.Fatal(err)
	}
	t := template.New("")
	t.Parse(templateLayout)
	t.Parse(templateContent)
	err = t.ExecuteTemplate(w, "layout", apv)

}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	type notFoundPageVars struct {
		globalPageVars
	}

	nfpv := notFoundPageVars{
		globalPageVars: gpv,
	}

	templateLayout, err := templateBox.FindString("layout.html")
	if err != nil {
		log.Fatal(err)
	}
	templateContent, err := templateBox.FindString("404.html")
	if err != nil {
		log.Fatal(err)
	}
	t := template.New("")
	t.Parse(templateLayout)
	t.Parse(templateContent)
	err = t.ExecuteTemplate(w, "layout", nfpv)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/posts/{slug}", postEndpoint).Methods("GET")
	router.HandleFunc("/about", aboutEndpoint).Methods("GET")
	router.HandleFunc("/", homeEndpoint).Methods("GET")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticBox)))

	router.PathPrefix("/").HandlerFunc(catchAllHandler)

	log.Println("Listening on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

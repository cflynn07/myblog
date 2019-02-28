package app

import (
	"bytes"
	"html/template"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

var postsBox = packr.New("Posts", "../posts")
var templateBox = packr.New("Templates", "../templates")

func truncHelper(s string) string {
	words := strings.Fields(s)
	maxPreviewLength := int(math.Min(40, float64(len(words))))
	words = words[0:maxPreviewLength]
	return strings.Join(words, " ")
}

// HomeHandler Handler for / route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
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

// PostHandler Handler for /posts/* route
func PostHandler(w http.ResponseWriter, r *http.Request) {
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

		t := template.New("")
		t.Parse(string(data))
		buf := bytes.NewBufferString("")
		err = t.Execute(buf, post)
		if err != nil {
			log.Fatal(err)
		}
		data = buf.Bytes()

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

// AboutHandler Handler for /about route
func AboutHandler(w http.ResponseWriter, r *http.Request) {
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

// CatchAllHandler Handler for undefined routes
func CatchAllHandler(w http.ResponseWriter, r *http.Request) {
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

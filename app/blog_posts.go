package app

import (
	"html/template"
	"time"
)

// Metadata for each blog post
type postData struct {
	Title          string
	Description    template.HTML
	Keywords       []string
	Date           time.Time
	Content        string
	ContentPreview string
}

// All published blog posts
type blogPosts map[string]*postData

// Blog uses this static map to display blog posts, the keys should match files
// in the template folder (sans the extension)
var bp = blogPosts{
	"2019-02-26-website-in-a-binary": &postData{
		Title:       "Bundling static website assets in a single binary with gobuffalo/packr",
		Description: template.HTML("How to use gobuffalo/packr to build a simple website in golang with static assets that's easy to test and can be bundled into a single binary for deployment."),
		Keywords:    []string{"golang", "go", "packr", "gobuffalo", "static assets"},
		Date:        time.Date(2019, time.February, 28, 0, 0, 0, 0, time.UTC),
	},
	"2019-03-04-helm-update-with-new-values": &postData{
		Title:       "HELM UPDATE WITH NEW VALUES",
		Description: template.HTML("How to use gobuffalo/packr to build a simple website in golang with static assets that's easy to test and can be bundled into a single binary for deployment."),
		Keywords:    []string{"golang", "", "static assets"},
		Date:        time.Date(2019, time.March, 5, 0, 0, 0, 0, time.UTC),
	},
}

var bpKeys = make([]string, len(bp))

func init() {
	bpKeys = []string{
		"2019-03-04-helm-update-with-new-values",
		"2019-02-26-website-in-a-binary",
	}
}

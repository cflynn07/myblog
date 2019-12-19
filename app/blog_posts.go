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
	"2019-09-10-circleci-to-github-actions": &postData{
		Title:       "CircleCI to Github Actions",
		Description: template.HTML(""),
		Keywords:    []string{""},
		Date:        time.Date(2019, time.September, 10, 0, 0, 0, 0, time.UTC),
	},
	"2019-08-30-automatically-initialized-and-version-controlled-database-in-kubernetes-and-helm-development-environments": &postData{
		Title:       "Automatically initalized and version controlled MySQL database in kubernetes and helm development environment",
		Description: template.HTML("A simple setup for an automatically initialized and version controlled MySQL database in kubernetes & helm development environments"),
		Keywords:    []string{"kubernetes", "helm", "skaffold", "development environment", "development environment mysql"},
		Date:        time.Date(2019, time.August, 30, 0, 0, 0, 0, time.UTC),
	},
	"2019-02-26-website-in-a-binary": &postData{
		Title:       "Bundling static website assets in a single binary with gobuffalo/packr",
		Description: template.HTML("How to use gobuffalo/packr to build a simple website in golang with static assets that's easy to test and can be bundled into a single binary for deployment."),
		Keywords:    []string{"golang", "go", "packr", "gobuffalo", "static assets"},
	},
	"2019-03-04-helm-update-with-new-values": &postData{
		Title:       "Deploying with helm from CI using --reuse-values and adding new values to your values.yaml file",
		Description: template.HTML("Adding new values to your values.yaml file and deploying from CI with --reuse-values can get you into trouble since tiller won't reference changes to values.yaml when creating kubernetes resources yaml."),
		Keywords:    []string{"kubernetes", "helm", "helm update", "CI deployment"},
		Date:        time.Date(2019, time.March, 5, 0, 0, 0, 0, time.UTC),
	},
}

var bpKeys = make([]string, len(bp))

func init() {
	bpKeys = []string{
		"2019-09-10-circleci-to-github-actions",
		"2019-08-30-automatically-initialized-and-version-controlled-database-in-kubernetes-and-helm-development-environments",
		"2019-03-04-helm-update-with-new-values",
		"2019-02-26-website-in-a-binary",
	}
}

package app

import (
	"html/template"
	"os"
)

// Template variables for all pages
type globalPageVars struct {
	Title           string
	Host            string
	Keywords        string
	Description     template.HTML
	GoogleAnalytics string
	DeploymentTime  string
	DeploymentSHA   string
}

var gpv = globalPageVars{
	Title:           "Casey Flynn",
	Host:            "https://cflynn.us",
	Keywords:        "Casey Flynn, blog, web development, programming, digital nomad",
	Description:     template.HTML("Casey Flynn's digital nomad programming, web development blog."),
	GoogleAnalytics: os.Getenv("GOOGLE_ANALYTICS"),
	DeploymentTime:  os.Getenv("DEPLOYMENT_TIME"),
	DeploymentSHA:   os.Getenv("DEPLOYMENT_SHA"),
}

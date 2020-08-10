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
	Image          string
}

type blogPosts map[string]*postData

var bpKeys = make([]string, 5)
var bp = make(blogPosts, 5)

func init() {
	post2020_07_20 := "2020-07-20-clubbingowl"
	post2020_07_14 := "2020-07-14-book-review-algorithms-in-a-nutshell"
	post2020_06_28 := "2020-06-28-book-review-the-go-programming-language"
	post2020_06_16 := "2020-06-16-book-review-learning-computer-architecture-with-raspberry-pi"
	post2020_04_26 := "2020-04-26-github-action-til-autoformat-readme"
	post2020_04_03 := "2020-04-03-book-review-fullstack-react"
	post2020_03_27 := "2020-03-27-quantifying-and-time-tracking-reading"
	post2020_03_14 := "2020-03-14-advanced-mysql-docker-tmux-demo"
	post2020_03_13 := "2020-03-13-book-review-high-performance-mysql"
	post2020_02_09 := "2020-02-09-book-review-mysql-crash-course"
	post2020_01_18 := "2020-01-18-new-tricks"
	post2019_12_23 := "2019-12-23-book-review-linux-command-line-shell-scripting-bible"
	post2019_09_10 := "2019-09-10-circleci-to-github-actions"
	post2019_08_30 := "2019-08-30-automatically-initialized-and-version-controlled-database-in-kubernetes-and-helm-development-environments"
	post2019_03_04 := "2019-03-04-helm-update-with-new-values"
	post2019_02_26 := "2019-02-26-website-in-a-binary"

	// Order determines post listing order. Newest entry at lowest index.
	bpKeys = []string{
		post2020_07_20,
		post2020_07_14,
		post2020_06_28,
		post2020_06_16,
		post2020_04_26,
		post2020_04_03,
		post2020_03_27,
		post2020_03_14,
		post2020_03_13,
		post2020_02_09,
		post2020_01_18,
		post2019_12_23,
		post2019_09_10,
		post2019_08_30,
		post2019_03_04,
		post2019_02_26,
	}

	bp[post2020_07_20] = &postData{
		Title:       "Starting my Software Engineer Career by Building an Complex, Technically Sophisticated and Over-engineered Product but a Terrible Business",
		Description: template.HTML(""),
		Keywords:    []string{""},
		Date:        time.Date(2020, time.July, 14, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/clubbingowl/clubbingowl_showcaseborder.png",
	}
	bp[post2020_07_14] = &postData{
		Title:       "Book Review: Algorithms in a Nutshell",
		Description: template.HTML("An \"In a Nutshell\" approach book might not be the best starting point for general learning"),
		Keywords:    []string{"Algorithms"},
		Date:        time.Date(2020, time.July, 14, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/algorithms_books.jpeg",
	}
	bp[post2020_06_28] = &postData{
		Title:       "Book Review: The Go Programming Language",
		Description: template.HTML("A comprehensive guide to understanding and using the Go programming language"),
		Keywords:    []string{"Go", "Golang", "Programming"},
		Date:        time.Date(2020, time.June, 28, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/2020-06-28/the_go_programming_language.jpg",
	}
	bp[post2020_06_16] = &postData{
		Title:       "Book Review: Learning Computer Architecture with Raspberry Pi",
		Description: template.HTML("A walkthrough of the history as well as the low level workings of computers"),
		Keywords:    []string{"Computer Architecture"},
		Date:        time.Date(2020, time.June, 16, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/2020-06-16/learning_computer_architecture.jpg",
	}
	bp[post2020_04_26] = &postData{
		Title:       "Creating a Reusable GitHub Action to Automatically Format a README for a TIL Repository",
		Description: template.HTML("GitHub Actions' ability to use docker containers can be exploited for many useful CI/CD tasks. This is an example of building an action that generates a formatted README for a repo of TILs"),
		Keywords:    []string{"GitHub", "GitHub Actions", "TIL", "Today I Learned"},
		Date:        time.Date(2020, time.April, 26, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/2020-04-26/Screen_Shot_2020-04-26_at_2.01.45_PM.png",
	}
	bp[post2020_04_03] = &postData{
		Title:       "Book Review: Fullstack React, The Complete Guide to ReactJS and Friends",
		Description: template.HTML("A guided tour of ReactJS philosophy and the modern web app frontend ecosystem"),
		Keywords:    []string{"Fullstack React", "ReactJS", "Book Review"},
		Date:        time.Date(2020, time.April, 4, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/fullstack-react_p.png",
	}
	bp[post2020_03_27] = &postData{
		Title:       "Quantifying and Time Tracking My Reading",
		Description: template.HTML("Using quantification and measuring techniques with Google Sheets to track my technical reading progress."),
		Keywords:    []string{},
		Date:        time.Date(2020, time.March, 27, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/fullstack-react-time-tracking_p.png",
	}
	bp[post2020_03_14] = &postData{
		Title:       "Advanced MySQL demo using docker, tmux, tmuxinator",
		Description: template.HTML("Using containers with tmux and tmuxinator to visualize demos"),
		Keywords:    []string{},
		Date:        time.Date(2020, time.March, 14, 0, 0, 0, 0, time.UTC),
	}
	bp[post2020_03_13] = &postData{
		Title:       "Book Review: High Performance MySQL and thoughts on digesting dense technical books",
		Description: template.HTML("High Performance MySQL by Baron Schwartz, Peter Zaitsev and Vadim Tkachenko. A deep dive into MySQL/RDBMSs and dense technical books"),
		Keywords:    []string{"book review"},
		Date:        time.Date(2020, time.March, 13, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/high-performance-mysql.png",
	}
	bp[post2020_02_09] = &postData{
		Title:       "Book Review: MySQL Crash Course",
		Description: template.HTML("My thoughts after reading MySQL Crash Course by Ben Forta"),
		Keywords:    []string{"book review"},
		Date:        time.Date(2020, time.February, 9, 0, 0, 0, 0, time.UTC),
	}
	bp[post2020_01_18] = &postData{
		Title:       "New Tricks",
		Description: template.HTML("A few new tricks and techniques I've recently incorporated into my workflow: (peco, yank, vim -, \"*yy register usage, hexyl, bat, bropages)"),
		Keywords:    []string{"shell scripting", "bash", "programming"},
		Date:        time.Date(2020, time.January, 18, 0, 0, 0, 0, time.UTC),
	}
	bp[post2019_12_23] = &postData{
		Title:       "Book Review: Linux Command Line and Shell Scripting Bible (3rd edition)",
		Description: template.HTML("My thoughts and insights from the technical book, \"Linux Command Line and Shell Scripting Bible (3rd edition)\""),
		Keywords:    []string{"book review", "shell scripting", "linux command line and shell scripting", "bash", "technical book review"},
		Date:        time.Date(2019, time.December, 23, 0, 0, 0, 0, time.UTC),
		Image:       "/static/images/linux-command-line-shell-scripting-bible.png",
	}
	bp[post2019_09_10] = &postData{
		Title:       "CircleCI to GitHub Actions",
		Description: template.HTML("Migrating a test, build and deploy pipeline from CircleCI to GitHub Actions"),
		Keywords:    []string{"GitHub Actions", "CircleCI", "CI", "CD"},
		Date:        time.Date(2019, time.September, 10, 0, 0, 0, 0, time.UTC),
	}
	bp[post2019_08_30] = &postData{
		Title:       "Automatically initalized and version controlled MySQL database in kubernetes+helm development environment",
		Description: template.HTML("A simple setup for an automatically initialized and version controlled MySQL database in kubernetes & helm development environments"),
		Keywords:    []string{"kubernetes", "helm", "skaffold", "development environment", "development environment mysql"},
		Date:        time.Date(2019, time.August, 30, 0, 0, 0, 0, time.UTC),
	}
	bp[post2019_03_04] = &postData{
		Title:       "Deploying to Kubernetes from CI with helm",
		Description: template.HTML("Adding new values to your values.yaml file and deploying from CI with --reuse-values can get you into trouble since tiller won't reference changes to values.yaml when creating kubernetes resources yaml"),
		Keywords:    []string{"kubernetes", "helm", "helm update", "CI deployment"},
		Date:        time.Date(2019, time.March, 5, 0, 0, 0, 0, time.UTC),
	}
	bp[post2019_02_26] = &postData{
		Title:       "Bundling static website assets in a single binary with gobuffalo/packr",
		Description: template.HTML("How to use gobuffalo/packr to build a simple website in golang with static assets that's easy to test and can be bundled into a single binary for deployment"),
		Keywords:    []string{"golang", "go", "packr", "gobuffalo", "static assets"},
		Date:        time.Date(2019, time.February, 26, 0, 0, 0, 0, time.UTC),
	}
}

package main

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

package app

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
	"2019-02-26-website-in-a-binary": &postData{
		Title:    "Website in a Binary",
		Subtitle: "subtitle subtitle subtitle",
		Keywords: "",
		Date:     "",
	},
}

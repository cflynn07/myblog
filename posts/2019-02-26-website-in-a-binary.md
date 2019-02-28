Most websites typically consist of server code and static files like css, js
and images. With golang, testing http route handlers that look for specific
files on disk, such as templates, can be a bit tricky. `go test` compiles your
go code and runs the binary from a temporary directory. It does not copy your
static files to the temporary directory, so route handlers wont be able to find
files on disk with relative include paths and your tests will fail. The same
problem exists with `go install`.

[gobuffalo/packr](https://github.com/gobuffalo/packr/) is a library that can
bundle static files inside binaries for production builds, and automatically
resolve include statments to absolute file paths so programs that depend on
files from the disk work no matter what context you run them in. You can even
use it with multi-stage docker builds to create a single binary containing all
your website's static assest in an alpine linux container.

The code in this project
[golang-packr-demo-fail](https://github.com/cflynn07/golang-packr-demo-fail)
uses the io/ioutils FileOpen() function in route handlers to load templates
with relative file paths.

<pre class="prettyprint mx-3 px-3 border-secondary rounded">
// handlers/handlers.go
package handlers

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// HomeHandler / route handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeVars := struct {
		Title string
	}{
		Title: "Home",
	}

  // relative include, only works if program is run from root of repository
	templateLayout, err := ioutil.ReadFile("templates/layout.html")
	if err != nil {
		log.Fatal(err)
	}
	templateHome, err := ioutil.ReadFile("templates/home.html")
	if err != nil {
		log.Fatal(err)
	}
	t := template.New("")
	t.Parse(string(templateLayout))
	t.Parse(string(templateHome))
	err = t.ExecuteTemplate(w, "layout", homeVars)
}
</pre>

This works if the program is executed within the
correct context, however it requires the templates and other static assets be
bundled with the built binary.

For example, the web server will work when `go run main.go` is run within the
repository, but if you run `go install` and try to run the server from your
PATH and the relative include paths wont resove, resulting in the server
crashing.

###### Working Example: `go run main.go` from within repository
<pre class="prettyprint mx-3 px-3 border-secondary rounded">
$ git clone git@github.com:cflynn07/golang-packr-demo-fail.git
$ cd golang-packr-demo-fail && go install
$ PORT=9000 go run main.go
2019/02/28 19:32:02 Listening port 9000
$ curl -X HEAD localhost:9000 -v
Warning: Setting custom HTTP method to HEAD with -X/--request may not work the
Warning: way you want. Consider using -I/--head instead.
* Rebuilt URL to: localhost:9000/
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 9000 (#0)
> HEAD / HTTP/1.1
> Host: localhost:9000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 28 Feb 2019 11:33:46 GMT
< Content-Length: 141
< Content-Type: text/html; charset=utf-8
</pre>

###### Failing Example: binary run from GOPATH/PATH
<pre class="prettyprint mx-3 px-3 border-secondary rounded">
$ git clone git@github.com:cflynn07/golang-packr-demo-fail.git
$ cd golang-packr-demo-fail && go install
$ go install.
$ PORT=9000 golang-packr-demo-fail
2019/02/28 19:37:21 Listening port 9000

// Then run
$ curl HEAD localhost:9000
curl: (7) Failed to connect to localhost port 9000: Connection refused

// And you will see your server failed with the following error
2019/02/28 19:37:23 open ../templates/layout.html: no such file or directory
</pre>


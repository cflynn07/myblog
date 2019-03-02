Most websites typically consist of server code and static files like css, js
and images. With golang, testing http route handlers that look for specific
files on disk, such as templates, can be a bit tricky. `go test` compiles your
go code and runs the binary from a temporary directory. It does not copy your
static files to the temporary directory, so route handlers wont be able to find
files on disk with relative include paths and your tests will fail. The same
problem exists with `go install`.

[gobuffalo/packr](https://github.com/gobuffalo/packr/) is a library that can
bundle static files inside binaries for production builds, but still load files
from disk during development, from any context, by automatically converting
relative include paths to absolute paths.

It's a great tool for building super lightweight images for deployments. With
multi-stage docker builds you can build your binary in a container that
contains the packr binary and the golang compiler, then copy the built binary
to an alpine linux image. The end product is an alpine linux container plus
just one additional file.

To illustrate some of the frustrations that can occur without packr, the code
in this project
[golang-packr-demo-fail](https://github.com/cflynn07/golang-packr-demo-fail)
uses the io/ioutils FileOpen() function in route handlers to load templates
with relative file paths. This file contains our route handler.

<pre class="prettyprint linenums">
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

This works if the program is executed within the correct context, however it
requires the templates and other static assets be bundled with the built
binary.

For example, the web server will work when `go run main.go` is run within the
repository, but if you run `go install` and try to run the server from your
PATH and the relative include paths wont resove, resulting in the server
crashing.

###### Working Example: `go run main.go` from within repository
<pre class="prettyprint lang-bsh">
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
<pre class="prettyprint lang-bsh">
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

#### packr is a solution to this

Packr exports a type `Box` which when initialized either computes absolute
paths to resources in a directory for later reference, or defers to the static
assets that have been included in the package via .go files if the `packr
build` step has been run.

This is convenient during development as files on disk will be continuously
reloaded every time they're requested by your golang code, and therefore will
reflect changes.

<pre class="prettyprint linenums">
// handlers/handlers.go
package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

// relative include path, converted to absolute path by packr
var templatesBox = packr.New("Templates", "../templates")

// HomeHandler / route handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeVars := struct {
		Title string
	}{
		Title: "Home",
	}

  // FindString either loads file from disk during development or 
  // from bundled .go file in production
	templateLayout, err := templatesBox.FindString("layout.html")
	if err != nil {
		log.Fatal(err)
	}
	templateHome, err := templatesBox.FindString("home.html")
	if err != nil {
		log.Fatal(err)
	}
	t := template.New("")
	t.Parse(templateLayout)
	t.Parse(templateHome)
	err = t.ExecuteTemplate(w, "layout", homeVars)
}
</pre>

To bundle your static assets with your golang code, first download the packr binary.
```
$ go get -u github.com/gobuffalo/packr/packr
```
Then from within your directory, run `packr build`.  

You'll see several new .go files, with `-packr.go` extensions. A quick peek inside these files shows the contents of your static files gzipped and stored in hexidecimal encodings.
<pre class="prettyprint linenums">
// +build !skippackr
// Code generated by github.com/gobuffalo/packr/v2. DO NOT EDIT.

// You can use the "packr2 clean" command to clean up this,
// and any other packr generated files.
package packrd

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file/resolver"
)

var _ = func() error {
	const gk = "458541e930fd91103b4fba34f99c4a9a"
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
		"2c8c91f95421405c8e06b6772791b143": "1f8b08000000000000ffaaae5648494dcbcc4b55504acecf2b49cd2b5152a8ade5b2c930b47bb277c1d3a57b6df4330cedb8aaab1552f352403280000000ffffa92b0fa131000000",
		"87b5eb07e6c9d4796e25fe54a00040f6": "1f8b08000000000000ff010000ffff0000000000000000",
		"d7c3cd8dceac76ab1c2129db9e0a26df": "1f8b08000000000000ff248e3f8fc2300c47f77c8adf79bf66bdc1ed727fd663280363680cad94a608cc5059f9ee2864b2f5fc2c3d3344b92c594029ecdb5309a538fef8f9ff1e4f875fccbaa6c1711d48215f7b924c154888830378150d98e6707f88f6741cff3ebfe87dd045930c66e8c6baa114f68d39f6ed9dcf5bdcab6c0695f596820a68dab24a5642875200c7be69ec5b8c1924c79af90a0000ffff28153539bd000000",
	})
	if err != nil {
		panic(err)
	}
	g.DefaultResolver = hgr

	func() {
		b := packr.New("Templates", "../templates")
		b.SetResolver("about.html", packr.Pointer{ForwardBox: gk, ForwardPath: "87b5eb07e6c9d4796e25fe54a00040f6"})
		b.SetResolver("home.html", packr.Pointer{ForwardBox: gk, ForwardPath: "2c8c91f95421405c8e06b6772791b143"})
		b.SetResolver("layout.html", packr.Pointer{ForwardBox: gk, ForwardPath: "d7c3cd8dceac76ab1c2129db9e0a26df"})
	}()

	return nil
}()
</pre>

At this point, if you run `go build` you'll have a binary that can be run from
anywhere and successfully serve your website and static files.

With multi-stage docker builds, you can produce really simple and minimal production docker images.

<pre class="prettyprint linenums">
FROM cflynnus/golang-1.12.0-packr2 as builder # Just golang:1.12.0 with the packr2 binary included
WORKDIR /go/src/golang-packr-demo
COPY . .
RUN packr2
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./main .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/golang-packr-demo/main ./main
CMD ["./main"]
</pre>

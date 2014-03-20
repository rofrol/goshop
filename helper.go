package main

import (
	"io"
	"net/http"
	"errors"
	"fmt"
	"os"
)

func serveError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Internal Server Error "+err.Error())
}

func redirectHandler(path string) func(http.ResponseWriter, *http.Request) {
	// http://stackoverflow.com/a/9936937/588759
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)
	}
	// usage: http.HandleFunc("/1", redirectHandler("/one"))
}

func parseForm(r *http.Request) error {
	// https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/73bqDlejJCQ
	// parseForm calls Request.ParseForm() excluding values from the URL query.
	// It returns an error if the request body was already parsed or if it failed
	// to parse.
	//
	// http://code.google.com/p/go/issues/detail?id=3630#c2
	// If I'm not mistaken, this "exploit" requires controlling the form's action.  If an attacker can control that, they could also probably redirect the user to their own server and steal all of the information and then redirect them back to the original action with properly-formed (but compromised) POST data.  If you are concerned about this in your webapps, it is probably trivial to add a quick `if r.Method = "POST" { r.URL.RawQuery = "" }`, though I would personally recommend auditing where the form tags get their action (in my own apps, it's always hard-coded in the template).
	// I think it's poor design to care where you get your form values.  I wouldn't mind if FormValues only got the data from the canonical source for the current method (GET -> query params, POST -> form body), but putting that in your code doesn't seem like the correct approach.  The PHP language (from my view) encourages people to care about the difference, but as soon as you do you make it harder to do simple things like unit test your code.  Often it is super easy to control form responses in query parameters for testing and they also are very good for creating links to pre-populate a form (akin to "mailto" links that provide the subject for you).  When you start caring where the data came from, the logic here becomes much more difficult.
	// Assuming that an established, authenticated and secure connection's $_POST could be trusted bit me once.... Never again.
	//
	// https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/ke_JP5IkofA
	// In the ParseForm method the values in the url query are overwritten by any values submitted via post
	// Nothing gets overwritten... Both values are added to req.Form. So with req.Form.Get you get the first value associated with the key - the one from your url query. The value from the post form is number two in the slice: req.Form["user"][1]
	if r.Form != nil {
		return errors.New("Request body was parsed already.")
	}
	tmp := r.URL.RawQuery
	r.URL.RawQuery = ""
	if err := r.ParseForm(); err != nil {
		return err
	}
	r.URL.RawQuery = tmp
	return nil
}

func serve404(w http.ResponseWriter) {
	// https://gist.github.com/1075842
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Nie ma takiej strony!")
}

func notlsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in notlsHandler")
	fullUrl := "https://localhost" + r.RequestURI
	http.Redirect(w, r, fullUrl, http.StatusMovedPermanently)
}

// http://stackoverflow.com/questions/13302020/rendering-css-in-a-go-web-application
type JustFilesFilesystem struct {
	fs http.FileSystem
}

func (fs JustFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

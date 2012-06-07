package main

import (
	"code.google.com/p/gorilla/mux"
	"code.google.com/p/gorilla/sessions"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
	"net/http"

//	"reflect"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func MyHandler(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "session-name")
	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it.
	session.Save(r, w)
}

func products(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	defer db.Close()

	sql := "select title, text, price from products order by title"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var title, text, price string
	for rows.Next() {
		rows.Scan(&title, &text, &price)
		levels = append(levels, map[string]string{"title": title, "text": text, "price": price})
	}

	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/products.html", "tpl/header.html", "tpl/footer.html")

	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	pageTemplate.Execute(w, map[string]interface{}{"levels": levels, "Header": "Products", "Copyright": "Roman Frołow"})

	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

func users(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	defer db.Close()

	sql := "select name1, surname from users order by surname"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var name1, surname string
	for rows.Next() {
		rows.Scan(&name1, &surname)
		levels = append(levels, map[string]string{"name1": name1, "surname": surname})
	}
	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/users.html", "tpl/header.html", "tpl/footer.html")

	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	pageTemplate.Execute(w, map[string]interface{}{"levels": levels, "Header": "Users", "Copyright": "Roman Frołow"})

	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

func productsAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	/*if not session.get('logged_in'):
	      abort(401)
	  g.db.execute('insert into products (title, text, price) values (?, ?, ?)',
	               [request.form['title'], request.form['text'], request.form['price']])
	  g.db.commit()
	  flash('New product was successfully added')
	  return redirect(url_for('show_products'))
	*/
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}

	stmt, err := tx.Prepare("insert into products(title, text, price) values(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	defer stmt.Close()
	// TODO
	/*
		if err := parseForm(r); err != nil {
			serveError(w, err)
			return
		}
	*/
	r.ParseForm()
	title := r.Form.Get("title")
	text := r.Form.Get("text")
	price := r.Form.Get("price")

	res, err := stmt.Exec(title, text, price)
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	last, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	fmt.Println("last", last)
	tx.Commit()

	session, _ := store.Get(r, "session-name")
	session.Values["last_product"] = last
	session.Save(r, w)

	// http://en.wikipedia.org/wiki/Post/Redirect/Get
	// http://en.wikipedia.org/wiki/HTTP_303
	// http://stackoverflow.com/questions/46582/response-redirect-with-post-instead-of-get
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "/" {
		serve404(w)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	session, _ := store.Get(r, "session-name")
	last := session.Values["last_product"]
	fmt.Println("last", last)

	pageTemplate, err := template.ParseFiles("tpl/index.html", "tpl/header.html", "tpl/footer.html")
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	pageTemplate.Execute(w, map[string]interface{}{"Header": "Home", "Copyright": "Roman Frołow"})
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

func hi(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hi")
}
func post(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	r.ParseForm()
	v := r.Form.Get("s")
	// If content-type not set manually, it will be guessed by http://golang.org/src/pkg/net/http/sniff.go
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, v)
	io.WriteString(w, `<br><form method="POST" action="/post"><input name="s"></form>`)
	io.WriteString(w, `<div>hello</div>`)
}

func serve404(w http.ResponseWriter) {
	// https://gist.github.com/1075842
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Nie ma takiej strony!")
}

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

var addr = flag.String("addr", ":9999", "http service address") // Q=17, R=18

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/", http.HandlerFunc(index))
	r.HandleFunc("/hi", http.HandlerFunc(hi))
	r.HandleFunc("/products", http.HandlerFunc(products))
	r.HandleFunc("/products/add", http.HandlerFunc(productsAdd))
	r.HandleFunc("/users", http.HandlerFunc(users))
	r.HandleFunc("/post", http.HandlerFunc(post))
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

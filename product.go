package main

import(
	"fmt"
	"html"
	"html/template"
	"net/http"
	"database/sql"
	"log"
)

func products(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Products", "Copyright": "Roman Frołow"}
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
	tplValues["levels"] = levels

	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/products.html", "tpl/header.html", "tpl/footer.html")
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	if _, ok := session.Values["login"]; ok {
		tplValues["login"] = session.Values["login"]
	}

	err = pageTemplate.Execute(w, tplValues)
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

	if err := r.ParseForm(); err != nil {
		serveError(w, err)
		return
	}
	title := html.EscapeString(r.Form.Get("title"))
	text := html.EscapeString(r.Form.Get("text"))
	price := html.EscapeString(r.Form.Get("price"))

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
	// 303 for HTTP 1.1, maybe problem with old corporate proxies, so 302 could be better
	//
	// https://groups.google.com/forum/?fromgroups#!msg/golang-nuts/HeAoybScSTU/qxp1H7mWZVYJ
	// The common practice is to redirect only after successful forms.
	// So forms with errors are treated by the same POST request, and so have
	// access to the data.
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}


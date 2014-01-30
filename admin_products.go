package main

import (
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html"
	"html/template"
	"log"
)

func admin_products(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Products", "Copyright": "Roman Frołow"}
	authorized := false
	if i, ok := session.Values["admin_login"]; ok {
		if i == "admin" {
			authorized = true
		}
		tplValues["admin_login"] = i
	}

	if ! authorized {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	defer db.Close()

	sql := "select title, text, price, quantity from products order by title"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var title, text, price, quantity string
	for rows.Next() {
		rows.Scan(&title, &text, &price, &quantity)
		levels = append(levels, map[string]string{"title": title, "text": text, "price": price, "quantity": quantity})
	}
	tplValues["levels"] = levels

	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/admin_products.html", "tpl/header.html", "tpl/admin_bar.html", "tpl/footer.html")
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	if i, ok := session.Values["admin_login"]; ok {
		tplValues["admin_login"] = i
	}

	err = pageTemplate.Execute(w, tplValues)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

func admin_products_add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	/*if not session.get('logged_in'):
	      abort(401)
	  g.db.execute('insert into products (title, text, price) values (?, ?, ?)',
	               [request.form['title'], request.form['text'], request.form['price']])
	  g.db.commit()
	  flash('New product was successfully added')
	  return redirect(url_for('show_products'))
	*/
	db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
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

	stmt, err := tx.Prepare("insert into products(title, text, price, quantity) values(?, ?, ?, ?)")
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
	quantity := html.EscapeString(r.Form.Get("quantity"))
	res, err := stmt.Exec(title, text, price, quantity)
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
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}


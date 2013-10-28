package main

import(
	"fmt"
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

	if i, ok := session.Values["login"]; ok {
		tplValues["login"] = i
	}

	err = pageTemplate.Execute(w, tplValues)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

package main

import (
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
)

func admin_orders(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Orders", "Copyright": "Roman Frołow"}
	authorized := false
	if i, ok := session.Values["admin_login"]; ok {
		if i == "admin" {
			authorized = true
		}
		tplValues["admin_login"] = i
	}

	if ! authorized {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}

	db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	defer db.Close()

	sql := "select user_id, product_id, price, sell_date from orders order by sell_date"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var user_id, product_id, price, sell_date string
	for rows.Next() {
		rows.Scan(&user_id, &product_id, &price, &sell_date)
		levels = append(levels, map[string]string{"user_id": user_id, "product_id": product_id, "price": price, "sell_date": sell_date})
	}
	tplValues["levels"] = levels

	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/admin_orders.html", "tpl/header.html", "tpl/admin_bar.html", "tpl/footer.html")
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

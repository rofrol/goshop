package main

import (
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
)

func admin_users(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Users", "Copyright": "Roman Frołow"}
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

	sql := "select login, password, name1, name2, surname from users order by login"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var login, password, name1, name2, surname string
	for rows.Next() {
		rows.Scan(&login, &password, &name1, &name2, &surname)
		levels = append(levels, map[string]string{"login": login, "password": password, "name1": name1, "name2": name2, "surname": surname})
	}
	tplValues["levels"] = levels
	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/admin_users.html", "tpl/admin_header.html", "tpl/footer.html")
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	pageTemplate.Execute(w, tplValues)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

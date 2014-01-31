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

	if err := r.ParseForm(); err != nil {
		serveError(w, err)
		return
	}

	fVals := map[string]string{}
	fVals["sent"] = html.EscapeString(r.Form.Get("sent"))
	if fVals["sent"] == "yes" {

		fVals["title"] = html.EscapeString(r.Form.Get("title"))
		fVals["description"] = html.EscapeString(r.Form.Get("description"))
		fVals["price"] = html.EscapeString(r.Form.Get("price"))
		fVals["quantity"] = html.EscapeString(r.Form.Get("quantity"))

		message_error := ""
		if fVals["title"] == "" {
			message_error += "\n<p>title can't be empty</p>"
		}
		if fVals["description"] == "" {
			message_error += "\n<p>description can't be empty</p>"
		}
		if fVals["price"] == "" {
			message_error += "\n<p>price can't be empty</p>"
		}
		if fVals["quantity"] == "" {
			message_error += "\n<p>quantity can't be empty</p>"
		}

		if message_error != "" {
			tplValues["message_error"] = template.HTML(message_error)
			tplValues["fVals"] = fVals
		} else {

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

			stmt, err := tx.Prepare("insert into products(title, description, price, quantity) values(?, ?, ?, ?)")
			if err != nil {
				fmt.Println(err)
				serveError(w, err)
				return
			}

			defer stmt.Close()

			res, err := stmt.Exec(fVals["title"], fVals["description"], fVals["price"], fVals["quantity"])
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

			session.Values["last_product"] = last
			session.Save(r, w)
			tplValues["message_info"] = template.HTML("<p>Adding succeeded</p>")
		}
	}

	levels, err := getProducts()
	tplValues["levels"] = levels
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	pageTemplate, err := template.ParseFiles("tpl/admin_products.html", "tpl/header.html", "tpl/admin_bar.html", "tpl/footer.html")
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	err = pageTemplate.Execute(w, tplValues)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

func getProducts() ([]map[string]string, error) {
	db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	sql := "select title, description, price, quantity from products order by title"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		return nil, err
	}
	defer rows.Close()

	levels := []map[string]string{}
	var title, description, price, quantity string
	for rows.Next() {
		rows.Scan(&title, &description, &price, &quantity)
		levels = append(levels, map[string]string{"title": title, "description": description, "price": price, "quantity": quantity})
	}
	return levels, nil
}

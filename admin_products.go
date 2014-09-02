package main

import (
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html"
	"html/template"
	"log"
	"io/ioutil"
	"io"
	"os"
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

	fVals := map[string]string{}
	if r.Method == "POST" {
		if false {
			hah, err := ioutil.ReadAll(r.Body);
			if err != nil {
				fmt.Printf("%s", err)
			}
			fmt.Printf("%v", string(hah))
			return
		}

		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			if part.FileName() == "" { // normal text field
				b, err := ioutil.ReadAll(part)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fVals[part.FormName()] = html.EscapeString(string(b))
			} else {
				//TODO: generate file name and store it with original name
				fileName := part.FileName()
				dst, err := os.Create("static/assets/" + fileName)
				defer dst.Close()

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if _, err := io.Copy(dst, part); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fVals[part.FormName()] = fileName
			}
		}
	}

	if fVals["sent"] == "yes" {

		message_error := []string{}
		if fVals["title"] == "" {
			message_error = append(message_error, "title can't be empty")
		}
		if fVals["description"] == "" {
			message_error = append(message_error, "description can't be empty")
		}
		if fVals["price"] == "" {
			message_error = append(message_error, "price can't be empty")
		}
		if fVals["quantity"] == "" {
			message_error = append(message_error, "quantity can't be empty")
		}
		if fVals["filename"] == "" {
			message_error = append(message_error, "filename can't be empty")
		}

		if len(message_error) != 0 {
			tplValues["message_error"] = message_error
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

			stmt, err := tx.Prepare("insert into products(title, description, price, quantity, filename) values(?, ?, ?, ?, ?)")
			if err != nil {
				fmt.Println(err)
				serveError(w, err)
				return
			}

			defer stmt.Close()

			res, err := stmt.Exec(fVals["title"], fVals["description"], fVals["price"], fVals["quantity"], fVals["filename"])
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
			tplValues["message_info"] = []string{"Adding succeeded"}
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

	sql := "select title, description, price, quantity, filename from products order by title"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		return nil, err
	}
	defer rows.Close()

	levels := []map[string]string{}
	var title, description, price, quantity, filename string
	for rows.Next() {
		rows.Scan(&title, &description, &price, &quantity, &filename)
		levels = append(levels, map[string]string{"title": title, "description": description, "price": price, "quantity": quantity, "filename": filename})
	}
	return levels, nil
}

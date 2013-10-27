package main

import (
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/sessions"
	"html"
	"html/template"
	"log"
)
func loginAvailable(login string) bool {
	if login == "" {
		fmt.Println("Pusty login.")
		return false
	}

	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("select count(login) from users where login = ?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(login)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}

	if count == 0 {
		return true
	}
	return false
}

func logged(r *http.Request, store *sessions.CookieStore) bool {
	session, _ := store.Get(r, "session-name")
	if _, ok := session.Values["login"]; ok {
		return true
	}
	return false
}

func params(r *http.Request, keys ...string) map[string]string {
	req := map[string]string{}
	for _, key := range keys {
		req[key] = html.EscapeString(r.Form.Get(key))
	}
	return req
}

func regParamsValid(req map[string]string) bool {
	return req["password"] != "" && req["password"] == req["repassword"] && req["login"] != "admin" && loginAvailable(req["login"])
}

func register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		serveError(w, err)
		return
	}
	//fmt.Println("Form.Values:", r.Form)
	session, _ := store.Get(r, "session-name")

	if logged(r, store) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if req := params(r, "login", "password", "repassword", "name1", "name2", "surname"); regParamsValid(req) {

		//if password good enough

		db, err := sql.Open("sqlite3", "./db/app.db")
		if err != nil {
			fmt.Println(err)
			serveError(w, err)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare("insert into users (login,password,name1,name2,surname) values (?,?,?,?,?)")
		if err != nil {
			fmt.Println(err)
			serveError(w, err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(req["login"], req["password"], req["name1"], req["name2"], req["surname"])
		if err != nil {
			fmt.Println(err)
			serveError(w, err)
			return
		}
		//last, err := res.LastInsertId()
		session.Values["req"] = req
		fmt.Println("req:")
		fmt.Println(session.Values["req"])
		/*
			req2 := session.Values["req"]
			login := req2.(map[string]interface{})["login"]
			fmt.Println("login:")
			fmt.Println(login)
		*/
		session.Save(r, w) // run before redirect
		http.Redirect(w, r, "/registered", http.StatusSeeOther)
	} else {
		pageTemplate, err := template.ParseFiles("tpl/register.html", "tpl/header.html", "tpl/footer.html")
		if err != nil {
			log.Fatalf("execution failed: %s", err)
			serveError(w, err)
		}

		tplValues := map[string]interface{}{"Header": "Register", "Copyright": "Roman Frołow"}
		if _, ok := session.Values["login"]; ok {
			tplValues["login"] = session.Values["login"]
		}

		pageTemplate.Execute(w, tplValues)
		if err != nil {
			log.Fatalf("execution failed: %s", err)
			serveError(w, err)
		}
	}
}

func registered(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	pageTemplate, err := template.ParseFiles("tpl/registered.html", "tpl/header.html", "tpl/footer.html")
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
	//req2 := session.Values["req"]
	//fmt.Println("login:")
	//fmt.Println(req2.(map[interface{}]string)["login"])

	tplValues := map[string]interface{}{"Header": "Registered", "Copyright": "Roman Frołow"}
	if req, ok := session.Values["req"]; ok {
		fmt.Println("req:")
		fmt.Println(req)
		fmt.Println("login:")
		fmt.Println(req.(map[string]string)["login"])
		tplValues["login"] = req.(map[string]string)["login"]
	}

	pageTemplate.Execute(w, tplValues)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	fmt.Println("session values:", session.Values)
	if err := r.ParseForm(); err != nil {
		serveError(w, err)
		return
	}
	login := html.EscapeString(r.Form.Get("login"))
	password := html.EscapeString(r.Form.Get("password"))

	if _, ok := session.Values["login"]; ok {
		fmt.Println("zalogowany")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if auth(login, password) {
		session.Values["login"] = login
		fmt.Println("session values:", session.Values)
		session.Save(r, w) // run before redirect
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		pageTemplate, err := template.ParseFiles("tpl/login.html", "tpl/header.html", "tpl/footer.html")
		if err != nil {
			log.Fatalf("execution failed: %s", err)
			serveError(w, err)
		}

		tplValues := map[string]interface{}{"Header": "Login", "Copyright": "Roman Frołow"}
		if _, ok := session.Values["login"]; ok {
			tplValues["login"] = session.Values["login"]
		}

		pageTemplate.Execute(w, tplValues)
		if err != nil {
			log.Fatalf("execution failed: %s", err)
			serveError(w, err)
		}
	}

}

func auth(login string, password string) bool {
	if login == "" || password == "" {
		fmt.Println("Pusty login i/lub hasło.")
		return false
	}

	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("select count(login) from users where login = ? and password = ?")
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer stmt.Close()

	rows, err := stmt.Query(login, password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var count int
	for rows.Next() {
		rows.Scan(&count)
	}

	if count == 1 {
		return true
	}
	return false
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "login")
	session.Save(r, w) // run before redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func users(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Users", "Copyright": "Roman Frołow"}
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		fmt.Println(err)
		serveError(w, err)
		return
	}
	defer db.Close()

	sql := "select name1, name2, surname from users order by surname"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var name1, name2, surname string
	for rows.Next() {
		rows.Scan(&name1, &name2, &surname)
		levels = append(levels, map[string]string{"name1": name1, "name2": name2, "surname": surname})
	}
	tplValues["levels"] = levels
	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/users.html", "tpl/header.html", "tpl/footer.html")
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}

	if _, ok := session.Values["login"]; ok {
		tplValues["login"] = session.Values["login"]
	}

	pageTemplate.Execute(w, tplValues)
	if err != nil {
		log.Fatalf("execution failed: %s", err)
		serveError(w, err)
	}
}


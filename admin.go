package main

import (
	"net/http"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html"
	"html/template"
	"log"
)

func admin_login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	fmt.Println("session values:", session.Values)
	if err := r.ParseForm(); err != nil {
		serveError(w, err)
		return
	}
	admin_login := html.EscapeString(r.Form.Get("admin_login"))
	fmt.Println("admin_login")
	fmt.Println(admin_login)
	password := html.EscapeString(r.Form.Get("password"))
	fmt.Println("password")
	fmt.Println(password)

	if _, ok := session.Values["admin_login"]; ok {
		fmt.Println("zalogowany")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else if auth(admin_login, password) {
		fmt.Println("auth")
		session.Values["admin_login"] = admin_login
		fmt.Println("session values:", session.Values)
		session.Save(r, w) // run before redirect
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		pageTemplate, err := template.ParseFiles("tpl/admin_login.html", "tpl/header.html", "tpl/admin_bar_notlogged.html", "tpl/footer.html")
		if err != nil {
			log.Fatalf("execution failed: %s", err)
			serveError(w, err)
		}

		tplValues := map[string]interface{}{"Header": "Admin", "Copyright": "Roman Frołow"}
		if _, ok := session.Values["login"]; ok {
			tplValues["admin_login"] = session.Values["admin_login"]
		}

		pageTemplate.Execute(w, tplValues)
		if err != nil {
			log.Fatalf("execution failed: %s", err)
			serveError(w, err)
		}
	}

}

func admin_logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "admin_login")
	session.Save(r, w) // run before redirect
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func admin_index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Admin Home", "Copyright": "Roman Frołow"}
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

	pageTemplate, err := template.ParseFiles("tpl/admin_index.html", "tpl/header.html", "tpl/admin_bar.html", "tpl/footer.html")
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

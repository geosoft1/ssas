package main

import (
	"net/http"
)

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templ.ExecuteTemplate(w, "signup", nil)
	case "POST":
		var user User
		user.Name = r.FormValue("name")
		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")
		if err := sqlInsert(&user); err != nil {
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("User created, <a href=\"/\">sign in</a>"))
	}
}

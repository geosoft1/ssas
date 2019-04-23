package main

import (
	"net/http"
)

func signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err, n := sqlUserCount(); err != nil || n == 0 {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		templ.ExecuteTemplate(w, "signin", nil)
	case "POST":
		var user User
		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")
		if err := sqlAuthenticateUser(&user); err != nil || !user.isActive {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		session, _ := store.Get(r, SESSID)
		session.Values["user"] = user
		session.Values["authenticated"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/a/home", http.StatusSeeOther)
	}
}

func signout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, SESSID)
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package main

import (
	"net/http"
)

func user(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, SESSID)
	user := session.Values["user"].(User)
	switch r.Method {
	case "GET":
		templ.ExecuteTemplate(w, "user", user)
	case "POST":
		switch r.FormValue("submit") {
		case "save":
			user.Name = r.FormValue("name")
			user.Password = r.FormValue("password")
			if err := sqlUpdateUser(&user); err != nil {
				http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
				return
			}
		case "delete":
			sqlDeleteUser(&user)
		}
		signout(w, r)
	}
}

package main

import (
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, SESSID)
	user := session.Values["user"].(User)
	templ.ExecuteTemplate(w, "home", user)
}

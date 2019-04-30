// +build random

package main

import (
	"log"
	"net/http"

	"github.com/geosoft1/token"

	gomail "gopkg.in/gomail.v2"
)

func reset(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templ.ExecuteTemplate(w, "reset", nil)
	case "POST":
		var user User
		user.Email = r.FormValue("email")
		if err := sqlGetUser(&user); err != nil {
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		// TODO reset password to random, uncomment the following 4 lines
		user.Password = token.GetToken(token_len)
		if err := sqlUpdateUser(&user); err != nil {
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}
		log.Println(user.Password)
		return
		// https://stackoverflow.com/a/24431749
		mail := gomail.NewMessage()
		mail.SetAddressHeader("From", config.SMTP.User, config.SMTP.Name)
		mail.SetAddressHeader("To", r.FormValue("email"), "")
		mail.SetHeader("Subject", "Check your password request")
		mail.SetBody("text/html", user.Password)
		dialer := gomail.NewPlainDialer(config.SMTP.Server, config.SMTP.Port, config.SMTP.User, config.SMTP.Password)
		//dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		if err := dialer.DialAndSend(mail); err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("A mail with instructions was send, read and <a href=\"/\">sign in</a>"))
	}
}

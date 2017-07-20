package main

import (
	"strings"
	"net/http"
)

type Page struct {
	Mess HeaderMessage
	Nav bool
	Info interface{}
}

type HeaderMessage struct {
	Visible string
	Type string
	Message string
}

func checkSession(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/signin", 302)		
	}
}

func signin(w http.ResponseWriter, r *http.Request) {
	p := &Page{HeaderMessage{Visible: "hidden"}, false, nil}
	if isSession(r) {
		http.Redirect(w, r, "/home", 302)
	}
	r.ParseForm()
	if checkPost(r.PostForm, "username", "password") {
	user := strings.TrimSpace(r.PostFormValue("username"))
		pass := strings.TrimSpace(r.PostFormValue("password"))
		if checkAccount(user, pass) {
			addSession(w)
			http.Redirect(w, r, "/home", 302)
		}
		p.Mess.Type = "Warning"
		p.Mess.Message = "Wrong Username/Password"
		p.Mess.Visible = ""
	}
}

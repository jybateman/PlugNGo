package main

import (
	"strings"
	"net/http"
)

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
	if !hasAdmin() {
		p.Mess.Type = "Warning"
		p.Mess.Message = "Please create an admin account"
		p.Mess.Visible = ""
		signupPage(w, r, p)
		return
	}
	r.ParseForm()
	if !checkPost(r.PostForm, "username", "password") {
		signinPage(w, r, p)
		return
	}
	user := strings.TrimSpace(r.PostFormValue("username"))
	pass := strings.TrimSpace(r.PostFormValue("password"))
	if checkAccount(user, pass) {
		addSession(w)
		http.Redirect(w, r, "/servers", 302)
	}
	p.Mess.Type = "Warning"
	p.Mess.Message = "Wrong Username/Password"
	p.Mess.Visible = ""
	signinPage(w, r, p)
}


func signup(w http.ResponseWriter, r *http.Request) {
	p := &Page{HeaderMessage{Visible: "hidden"}, false, nil}
	if isSession(r) {
		http.Redirect(w, r, "/home", 302)
	}
	if !hasAdmin() {
		p.Mess.Type = "Warning"
		p.Mess.Message = "Admin account already exist"
		p.Mess.Visible = ""
		signinPage(w, r, p)
		return
	}
	r.ParseForm()
	if !checkPost(r.PostForm, "newusername", "password", "confpassword") {
		signupPage(w, r, p)
		return
	}
	user := strings.TrimSpace(r.PostFormValue("newusername"))
	pass := strings.TrimSpace(r.PostFormValue("password"))
	if err := addAdmin(user, pass); err == nil {
		addSession(w)
		http.Redirect(w, r, "/home", 302)
	}
	p.Mess.Type = "Danger"
	p.Mess.Message = "Couldn't create account"
	p.Mess.Visible = ""
	signupPage(w, r, p)
}

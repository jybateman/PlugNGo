package main

import (
	"log"
	"time"
	"strings"
	"net/http"
)

func checkSession(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/signin", 302)
	}
	http.Redirect(w, r, "/home", 302)
}

func signin(w http.ResponseWriter, r *http.Request) {
	p := &Page{HeaderMessage{Visible: "hidden"}, false, nil, nil}
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
		http.Redirect(w, r, "/home", 302)
	}
	p.Mess.Type = "Warning"
	p.Mess.Message = "Wrong Username/Password"
	p.Mess.Visible = ""
	signinPage(w, r, p)
}

func signup(w http.ResponseWriter, r *http.Request) {
	p := &Page{HeaderMessage{Visible: "hidden"}, false, nil, nil}
	if isSession(r) {
		http.Redirect(w, r, "/home", 302)
	}
	if hasAdmin() {
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
	pass := strings.TrimSpace(r.PostFormValue("password"))
	confpass := strings.TrimSpace(r.PostFormValue("confpassword"))
	if strings.Compare(pass, confpass) != 0 {
		p.Mess.Type = "Danger"
		p.Mess.Message = "Password does not match the confirm password"
		p.Mess.Visible = ""
		signupPage(w, r, p)
		return
	}
	user := strings.TrimSpace(r.PostFormValue("newusername"))
	if err := addAdmin(user, pass); err == nil {
		addSession(w)
		http.Redirect(w, r, "/home", 302)
	}
	p.Mess.Type = "Danger"
	p.Mess.Message = "Couldn't create account"
	p.Mess.Visible = ""
	signupPage(w, r, p)
}

func home(w http.ResponseWriter, r *http.Request) {
	p := &Page{HeaderMessage{Visible: "hidden"}, false, nil, nil}
	if !isSession(r) {
		p.Mess.Type = "Warning"
		p.Mess.Message = "Please signin using Admin account"
		p.Mess.Visible = ""
		signinPage(w, r, p)
		return
	}
	p.Nav = true
	p.Info = plugs
	homePage(w, r, p)
}

func plug(w http.ResponseWriter, r *http.Request) {
	p := &Page{HeaderMessage{Visible: "hidden"}, false, nil, nil}
	if !isSession(r) {
		p.Mess.Type = "Warning"
		p.Mess.Message = "Please signin using Admin account"
		p.Mess.Visible = ""
		signinPage(w, r, p)
		return
	}
	path := strings.SplitAfter(r.URL.Path, "/")
	log.Println(len(path), path[2])
	p.Nav = true
	if len(path) < 3 {
		p.Mess.Type = "Warning"
		p.Mess.Message = "Unknown plug ID"
		p.Mess.Visible = ""
		homePage(w, r, p)
		return
	}
	pl, ok := plugs[path[2]]
	if !ok {
		p.Mess.Type = "Danger"
		p.Mess.Message = "Failed to retrieve plug information"
		p.Mess.Visible = ""
		homePage(w, r, p)
		return
	}
	go pl.GetSchedule()
	time.Sleep(time.Millisecond*500)
	p.Info = plugs
	p.Extra = pl
	plugPage(w, r, p)
}

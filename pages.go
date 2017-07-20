package main

 import (
	"log"
	"net/http"
	"html/template"
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

func signinPage(w http.ResponseWriter, r *http.Request, p *Page) {
	tpl, err := template.ParseFiles("html/signin.html", "html/header.html")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	tpl.Execute(w, p)
}

func signupPage(w http.ResponseWriter, r *http.Request, p *Page) {
	tpl, err := template.ParseFiles("html/signup.html", "html/header.html")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	tpl.Execute(w, p)
}

func homePage(w http.ResponseWriter, r *http.Request, p *Page) {
	tpl, err := template.ParseFiles("html/home.html", "html/header.html")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	tpl.Execute(w, p)
}

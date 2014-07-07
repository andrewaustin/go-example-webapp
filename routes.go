package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func slash(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	name := params["name"]

	addName(name)

	w.Write([]byte("Hello, " + name))
}

func saidHelloTo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strings.Join(getNames(), "\n")))
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
	}{
		"BOB DOLE",
	}

	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("One day login..."))
}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("One day logout"))
}

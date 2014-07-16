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
	name := "World!"

	session, _ := store.Get(r, "user")
	id := session.Values["id"]

	if id != nil {
		tempName := getUsername(id.(int))
		if len(tempName) > 0 {
			name = tempName
		}
	}

	w.Write([]byte("Hello, " + name))
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
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if len(user) == 0 || len(pass) == 0 {
		http.Error(w, "Please supply non-empty username and password", 401)
	} else {
		// TODO: This is bad for security =)
		// User bcrypt later
		if pass == getUserPass(user) {
			session, err := store.Get(r, "user")
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Fatal(err)
			}

			session.Values["id"] = getUserId(user)
			session.Save(r, w)

			http.Redirect(w, r, "/", 302)

		} else {

			http.Error(w, "Invalid Username or Password", 401)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")
	delete(session.Values, "id")
	session.Save(r, w)
	w.Write([]byte("Logged Out"))
}

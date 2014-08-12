package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type appHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		log.Println(err)
		switch status {
		case http.StatusNotFound:
			http.Error(w, http.StatusText(http.StatusNotFound), 404)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
		}
	}
}

func slash(w http.ResponseWriter, r *http.Request) (int, error) {
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

	return http.StatusOK, nil
}

func hello(w http.ResponseWriter, r *http.Request) (int, error) {
	params := mux.Vars(r)

	name := params["name"]

	addName(name)

	w.Write([]byte("Hello, " + name))

	return http.StatusOK, nil
}

func saidHelloTo(w http.ResponseWriter, r *http.Request) (int, error) {
	w.Write([]byte(strings.Join(getNames(), "\n")))

	return http.StatusOK, nil
}

func loginPage(w http.ResponseWriter, r *http.Request) (int, error) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if err = t.Execute(w, nil); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if len(user) == 0 || len(pass) == 0 {
		http.Error(w, "Please supply non-empty username and password", 401)
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(getUserHash(user)), []byte(pass))
		if err == nil {
			session, err := store.Get(r, "user")
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Fatal(err)
			}

			session.Values["id"] = getUserId(user)
			session.Save(r, w)

			http.Redirect(w, r, "/", 302)

		} else {
			log.Println(err)
			http.Error(w, "Invalid Username or Password", 401)
		}
	}
}

func registerPage(w http.ResponseWriter, r *http.Request) (int, error) {
	t, err := template.ParseFiles("templates/register.html")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = t.Execute(w, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func register(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if len(user) == 0 || len(pass) == 0 {
		http.Error(w, "Please supply non-empty username and password", 401)
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
		if err == nil {
			setUserPassword([]byte(user), hashedPassword)
			http.Redirect(w, r, "/login", 302)
		} else {
			log.Fatal(err)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) (int, error) {
	session, _ := store.Get(r, "user")
	delete(session.Values, "id")
	session.Save(r, w)
	w.Write([]byte("Logged Out"))

	return http.StatusOK, nil
}

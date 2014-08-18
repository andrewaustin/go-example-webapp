package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type appHandler struct {
	*appContext
	handler func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.handler(ah.appContext, w, r)
	if err != nil {
		log.Println(err)
		switch status {
		case http.StatusNotFound:
			http.Error(w, http.StatusText(http.StatusNotFound), 404)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
		}
	}
}

func slash(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	name := "World!"

	session, _ := a.store.Get(r, "user")
	id := session.Values["id"]

	if id != nil {
		tempName, err := getUsername(a.db, id.(int))
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if len(tempName) > 0 {
			name = tempName
		}
	}

	w.Write([]byte("Hello, " + name))

	return http.StatusOK, nil
}

func hello(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	params := mux.Vars(r)

	name := params["name"]

	addName(a.db, name)

	w.Write([]byte("Hello, " + name))

	return http.StatusOK, nil
}

func saidHelloTo(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	names, err := getNames(a.db)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Write([]byte(strings.Join(names, "\n")))

	return http.StatusOK, nil
}

func loginPage(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if err = t.Execute(w, nil); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func login(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if len(user) == 0 || len(pass) == 0 {
		return http.StatusUnauthorized, errors.New("Please supply non-empty username and password")
	} else {
		userHash, err := getUserHash(a.db, user)
		if err != nil {
			return http.StatusUnauthorized, err
		}
		err = bcrypt.CompareHashAndPassword([]byte(userHash), []byte(pass))
		if err == nil {
			session, err := a.store.Get(r, "user")
			if err != nil {
				return http.StatusInternalServerError, err
			}

			var userId int
			userId, err = getUserId(a.db, user)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			session.Values["id"] = userId
			session.Save(r, w)

			http.Redirect(w, r, "/", 302)

		} else {
			return http.StatusUnauthorized, err
		}
	}

	return http.StatusOK, nil
}

func registerPage(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
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

func register(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if len(user) == 0 || len(pass) == 0 {
		return http.StatusUnauthorized, errors.New("Please supply non-empty username and password")
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
		if err == nil {
			err = setUserPassword(a.db, []byte(user), hashedPassword)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			http.Redirect(w, r, "/login", 302)
		} else {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func logout(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	session, _ := a.store.Get(r, "user")
	delete(session.Values, "id")
	session.Save(r, w)
	w.Write([]byte("Logged Out"))

	return http.StatusOK, nil
}

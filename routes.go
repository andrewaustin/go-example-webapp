package main

import (
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

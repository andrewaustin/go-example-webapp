package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var conf Config
var db *sql.DB
var connStr string
var store *sessions.CookieStore

type Config struct {
	Port         int
	User         string
	Pass         string
	Database     string
	CookieSecret string
}

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}

	connStr = conf.User + ":" + conf.Pass + "@tcp(127.0.0.1:3306)/" + conf.Database
	store = sessions.NewCookieStore([]byte(conf.CookieSecret))
}

func main() {
	var err error

	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Establishing database connection...")

	router := mux.NewRouter()
	router.Handle("/", appHandler(slash))
	router.Handle("/hello/{name:[a-zA-Z]+}", appHandler(hello))
	router.Handle("/who", appHandler(saidHelloTo))
	router.Handle("/login", appHandler(loginPage)).Methods("GET")
	router.Handle("/login", appHandler(login)).Methods("POST")
	router.Handle("/register", appHandler(registerPage)).Methods("GET")
	router.Handle("/register", appHandler(register)).Methods("POST")
	router.Handle("/logout", appHandler(logout))
	log.Println("Listening...")
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}

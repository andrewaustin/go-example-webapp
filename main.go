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

type appContext struct {
	conf    Config
	db      *sql.DB
	connStr string
	store   *sessions.CookieStore
}

func NewAppContext(conf Config) *appContext {
	connStr := conf.User + ":" + conf.Pass + "@tcp(127.0.0.1:3306)/" + conf.Database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	store := sessions.NewCookieStore([]byte(conf.CookieSecret))

	return &appContext{conf, db, connStr, store}
}

type Config struct {
	Port         int
	User         string
	Pass         string
	Database     string
	CookieSecret string
}

func main() {
	var err error
	var config Config

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}

	appCtx := NewAppContext(config)
	defer appCtx.db.Close()
	router := mux.NewRouter()
	router.Handle("/", appHandler{appCtx, slash})
	router.Handle("/hello/{name:[a-zA-Z]+}", appHandler{appCtx, hello})
	router.Handle("/who", appHandler{appCtx, saidHelloTo})
	router.Handle("/login", appHandler{appCtx, loginPage}).Methods("GET")
	router.Handle("/login", appHandler{appCtx, login}).Methods("POST")
	router.Handle("/register", appHandler{appCtx, registerPage}).Methods("GET")
	router.Handle("/register", appHandler{appCtx, register}).Methods("POST")
	router.Handle("/logout", appHandler{appCtx, logout})
	log.Println("Listening...")
	err = http.ListenAndServe(":"+strconv.Itoa(appCtx.conf.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}

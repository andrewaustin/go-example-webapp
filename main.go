package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var conf Config
var db *sql.DB
var connStr string

type Config struct {
	Port     int
	User     string
	Pass     string
	Database string
	DBSecret string
}

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Println(err)
	}

	connStr = conf.User + ":" + conf.Pass + "@tcp(127.0.0.1:3306)/" + conf.Database
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
	router.HandleFunc("/", slash)
	router.HandleFunc("/hello/{name:[a-zA-Z]+}", hello)
	router.HandleFunc("/who", saidHelloTo)
	router.HandleFunc("/login", loginPage)
	log.Println("Listening...")
	http.ListenAndServe(":"+strconv.Itoa(conf.Port), router)
}
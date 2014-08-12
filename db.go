package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func addName(name string) {
	stmt, err := db.Prepare("INSERT INTO hellos(name) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}
}

func getNames() []string {
	names := make([]string, 0)
	var name string

	rows, err := db.Query("SELECT name FROM hellos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return names
}

func getUserHash(user string) string {
	var pass string

	err := db.QueryRow("SELECT password FROM users WHERE username=?", user).Scan(&pass)
	if err != nil {
		log.Println(err)
	}

	return pass
}

func setUserPassword(user, hash []byte) {
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user, hash)
	if err != nil {
		log.Println("Error here...")
		log.Fatal(err)
	}
}

func getUserId(user string) int {
	var id int

	err := db.QueryRow("SELECT id FROM users WHERE username=?", user).Scan(&id)
	if err != nil {
		log.Println(err)
	}

	return id
}

func getUsername(id int) string {
	var name string

	err := db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&name)
	if err != nil {
		log.Println(err)
	}

	return name
}

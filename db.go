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

	rows, err := db.Query("select name from hellos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(name)
		names = append(names, name)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return names
}

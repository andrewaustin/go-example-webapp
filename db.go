package main

import _ "github.com/go-sql-driver/mysql"

func addName(name string) error {
	stmt, err := db.Prepare("INSERT INTO hellos(name) VALUES(?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

func getNames() ([]string, error) {
	names := make([]string, 0)
	var name string

	rows, err := db.Query("SELECT name FROM hellos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return names, nil
}

func getUserHash(user string) (string, error) {
	var pass string

	err := db.QueryRow("SELECT password FROM users WHERE username=?", user).Scan(&pass)
	if err != nil {
		return "", err
	}

	return pass, nil
}

func setUserPassword(user, hash []byte) error {
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES(?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user, hash)
	if err != nil {
		return err
	}

	return nil
}

func getUserId(user string) (int, error) {
	var id int

	err := db.QueryRow("SELECT id FROM users WHERE username=?", user).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func getUsername(id int) (string, error) {
	var name string

	err := db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

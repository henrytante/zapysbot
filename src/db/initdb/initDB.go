package initdb

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB()  {
	db, err := sql.Open("sqlite3", "./sistema.db")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	usersaldo := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		saldo int
	);
	`
	produtos := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type text,
		email text,
		pass text
	);
	`

	_, err = db.Exec(usersaldo)
	if err != nil{
		log.Fatal(err)
	}
	_, err = db.Exec(produtos)
	if err != nil{
		log.Fatal(err)
	}
}
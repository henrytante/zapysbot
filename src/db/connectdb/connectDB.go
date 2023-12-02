package connectdb

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./sistema.db")
	if err != nil{
		return nil, err
	}
	err = db.Ping()
	if err != nil{
		return nil, err
	}
	return db, nil
}
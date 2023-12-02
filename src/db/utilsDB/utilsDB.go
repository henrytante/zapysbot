package utilsdb

import (
	"fmt"
	"log"
	"telebotgo/src/db/connectdb"

	_ "github.com/mattn/go-sqlite3"
)


func GetUserSaldo(ID int64) (int64, error) {
	db, err := connectdb.ConnectDB()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	var saldo int64
	err = db.QueryRow("SELECT saldo FROM users where id = ?", ID).Scan(&saldo)
	if err != nil{
		return 0, err
	}
	return saldo, nil
}
func InsertNewUser(ID int64) error {
	var saldo int64
	
	db, err := connectdb.ConnectDB()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", ID). Scan(&count)
	if err != nil{
		log.Fatal(err)
	}
	if count == 0{
		if _, err = db.Exec("INSERT INTO users (id, saldo) values(?, ?)", ID, saldo); err != nil{
			return err
		}
		fmt.Println("Novo usuario inserido, ID:", ID)
	}
	
	return nil
}
func AddSaldo(ID int64) error {
	var saldo int64
	saldo = 10
	db, err := connectdb.ConnectDB()
	if err != nil{
		return err
	}
	defer db.Close()
	if _, err = db.Exec("UPDATE users SET saldo = saldo + ? WHERE id = ?", saldo, ID); err != nil{
		return err
	}
	return nil
}
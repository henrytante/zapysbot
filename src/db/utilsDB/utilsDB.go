package utilsdb

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"telebotgo/src/db/connectdb"
)

func GetUserSaldo(ID int64) (int64, error) {
	db, err := connectdb.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var saldo int64
	err = db.QueryRow("SELECT saldo FROM users where id = ?", ID).Scan(&saldo)
	if err != nil {
		return 0, err
	}
	return saldo, nil
}
func InsertNewUser(ID int64) error {
	var saldo int64

	db, err := connectdb.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", ID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		if _, err = db.Exec("INSERT INTO users (id, saldo) values(?, ?)", ID, saldo); err != nil {
			return err
		}
		fmt.Println("Novo usuario inserido, ID:", ID)
	}

	return nil
}
func AddSaldo(ID, valor int64) error {
	db, err := connectdb.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec("UPDATE users SET saldo = saldo + ? WHERE id = ?", valor, ID); err != nil {
		return err
	}
	return nil
}

func InsertProduct(email, senha, plataforma string, valor int) error {
	db, err := connectdb.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec("INSERT INTO products (plataforma, email, pass, preco) values (?,?,?,?)", plataforma, email, senha, valor); err != nil {
		return err
	}
	return nil
}
func DeleteProduct(id int) error {
	db, err := connectdb.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec("DELETE from products WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

type Login struct {
	ID         int
	Plataforma string
	Email      string
	Senha      string
}

func Logins() ([]Login, error) {
	var logins []Login

	db, err := connectdb.ConnectDB()
	if err != nil {
		return logins, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,plataforma, email, pass FROM products")
	if err != nil {
		return logins, err
	}
	defer rows.Close()

	for rows.Next() {
		var l Login
		err := rows.Scan(&l.ID, &l.Plataforma, &l.Email, &l.Senha)
		if err != nil {
			return logins, err
		}
		logins = append(logins, l)
	}

	if err = rows.Err(); err != nil {
		return logins, err
	}

	return logins, nil
}

type PlataformaInfo struct {
	Quantidade int
	Preco      int
}

func CountProductsByPlatform() (map[string]PlataformaInfo, error) {
	db, err := connectdb.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	countByPlatform := make(map[string]PlataformaInfo)

	query := `
        SELECT plataforma, COUNT(*) AS quantidade, AVG(preco) AS preco_medio
        FROM products
        GROUP BY plataforma
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var plataforma string
		var quantidade int
		var precoMedio float64

		if err := rows.Scan(&plataforma, &quantidade, &precoMedio); err != nil {
			return nil, err
		}

		// Converte o preço médio de float64 para int
		preco := int(precoMedio)

		// Atualiza o mapa com as informações da plataforma
		info := PlataformaInfo{Quantidade: quantidade, Preco: preco}
		countByPlatform[plataforma] = info
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return countByPlatform, nil
}

func ComprarLogin(plataforma string) (Login, error) {
	var login Login
	db, err := connectdb.ConnectDB()
	if err != nil {
		return login, err
	}
	defer db.Close()

	// Passo 1: Seleciona um registro aleatório
	row := db.QueryRow("SELECT email, pass FROM products WHERE plataforma = ? ORDER BY RANDOM() LIMIT 1", plataforma)
	err = row.Scan(&login.Email, &login.Senha)
	if err != nil {
		return login, err
	}

	// Passo 2: Deleta o registro baseado no email e senha obtidos
	_, err = db.Exec("DELETE FROM products WHERE email = ? AND pass = ?", login.Email, login.Senha)
	if err != nil {
		return login, err
	}

	return login, nil
}

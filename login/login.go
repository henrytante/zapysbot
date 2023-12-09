package login

import (
	"database/sql"
	"fmt"
	"log"
	"telebotgo/utils/clear"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
)

func Logar() {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	db, err := sql.Open("mysql", "")
	if err != nil {
		log.Fatal("Erro na conexão com o sistema de login")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao tentar dar vida à conexão com o banco de dados")
	}
	clear.Clear()
	fmt.Println(green("Sistema de login ON-LINE!\n"))
	var user, pass string
	maxTentativas := 3
	limiteTentativasExcedido := false
    for tentativas := 0; tentativas < maxTentativas; tentativas++ {
		fmt.Print("\n\nDigite seu usuário: ")
		fmt.Scan(&user)
		fmt.Print("Digite sua senha: ")
		fmt.Scan(&pass)

		var retrievedUser string
		err := db.QueryRow("SELECT user FROM paybotusers WHERE user = ? AND pass = ?", user, pass).Scan(&retrievedUser)
		switch {
		case err == sql.ErrNoRows:
			clear.Clear()
			fmt.Println(red("Login inválido. Tente novamente. Caso não tenha um, contato: https://t.me/h1000dev\n"))
			fmt.Printf(yellow("Tentativas restantes: ", maxTentativas-(tentativas+1)))
		case err != nil:
			log.Fatal(err)
		default:
			clear.Clear()
			fmt.Printf(green("Login bem-sucedido! Seja bem-vindo ", user))
			return // Encerra a função após o login bem-sucedido
		}

		if tentativas == maxTentativas-1 {
			limiteTentativasExcedido = true
			break // Sai do loop quando o limite de tentativas for atingido
		}
	}

	if limiteTentativasExcedido {
		clear.Clear()
		log.Fatal("Limite de tentativas excedido.")
	}
}


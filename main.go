package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	
	botcontroller "telebotgo/src/controller/botController"
	"telebotgo/src/db/initdb"
	utilsdb "telebotgo/src/db/utilsDB"
	"telebotgo/utils/clear"
	_ "telebotgo/utils/clear"

	"github.com/fatih/color"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	
	initdb.InitDB()
}

func main() {
	green := color.New(color.FgGreen).SprintFunc()
	var op int
	var email, senha, plataforma string
	var valor int
	var token string
	var mpToken string
	clear.Clear()
	fmt.Println("O que deseja fazer?")
	fmt.Println("1 - Iniciar bot\n2 - Logins")
	fmt.Print("-> ")
	fmt.Scan(&op)
	clear.Clear()
	switch op {
	case 1:
		fmt.Print("\nDigite seu token do telegram: ")
		fmt.Scan(&token)
		fmt.Print("\nDigite seu token do mercado pago: ")
		fmt.Scan(&mpToken)
		bot, err := tg.NewBotAPI(token)
		if err != nil {
			log.Fatal(err)
		}
		clear.Clear()
		log.Println(green("Autenticado com BOT: ", bot.Self.UserName))
		updateConfig := tg.NewUpdate(0)
		updateConfig.Timeout = 60
		updates, err := bot.GetUpdatesChan(updateConfig)
		if err != nil {
			log.Fatal(err)
		}
		for update := range updates {
			if update.CallbackQuery != nil {
				botcontroller.HandleCallbackQuery(update.CallbackQuery)
			} else if update.Message != nil {
				botcontroller.HandleMessage(update, bot, mpToken)
			}
		}
	case 2:
		for {

			fmt.Println("\nO que deseja fazer?")
			fmt.Println("1 - Adicionar login\n2 - Remover login\n3 - Ver logins\n0 - Sair")
			fmt.Print("-> ")
			var id int
			fmt.Scan(&op)
			switch op {
			case 1:
				fmt.Print("Digite a plataforma: ")
				fmt.Scan(&plataforma)
				fmt.Print("Digite o email: ")
				fmt.Scan(&email)
				fmt.Print("Digite a senha: ")
				fmt.Scan(&senha)
				fmt.Print("Digite o preço: ")
				fmt.Scan(&valor)

				if err := utilsdb.InsertProduct(email, senha, plataforma, valor); err != nil {
					fmt.Println("Erro ao adicionar login")
					panic(err)
				} else {
					clear.Clear()
					fmt.Println("Login cadastrado")
				}

			case 2:
				fmt.Print("Digite o id do login que deseja remover: ")
				fmt.Scan(&id)
				if err := utilsdb.DeleteProduct(id); err != nil {
					fmt.Println("Erro ao remover login")
					panic(err)
				} else {
					clear.Clear()
					if err == sql.ErrNoRows {
						fmt.Println("Login não encontrado")
					} else {
						fmt.Println("Login deletado")
					}
				}
			case 3:
				logins, err := utilsdb.Logins()
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("Nenhum login existente")
						return
					}
					log.Fatal("Erro ao ler logins:", err)
				}
				clear.Clear()
				fmt.Println("Logins: ")
				for _, l := range logins {
					fmt.Printf("ID: %d,Plataforma: %s, Email: %s, Senha: %s\n", l.ID,l.Plataforma, l.Email, l.Senha )
				}
			case 0:
				os.Exit(0)
			}
		}
	}
}

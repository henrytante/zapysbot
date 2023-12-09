package botcontroller

import (
	"fmt"
	"log"
	"telebotgo/src/db/connectdb"
	utilsdb "telebotgo/src/db/utilsDB"
	botview "telebotgo/src/view/bot_view"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func produtosArray(plataforma string, array []string) bool {
	for _, elemento := range array{
		if elemento == plataforma{
			return true
		}
	}
	return false
}

func SaldoSuficiente(saldoAtual, precoLogin int) bool {
    if saldoAtual < precoLogin{
		return false
	}else{
		return true
	}
}


func ComprarCommand(update tg.Update, bot *tg.BotAPI, userID int) {
	db, err := connectdb.ConnectDB()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	var plataformas []string
	var saldoAtual int
	func ()  {
		counts, err := utilsdb.CountProductsByPlatform()
			if err != nil {
				panic(err)
			}
			for plataforma, _ := range counts {
				plataformas = append(plataformas, plataforma)
			}
	}()
	comando := update.Message.CommandArguments()
	
	if len(comando) == 0 {
		botview.SendReply(bot, update.Message.Chat.ID, "Comandos de Compras💵\n\n/comprar produtos <- Esse comando mostra o estoque\n/comprar {login desejado} <- Compra um login\n\nSe não tiver saldo suficiente use /recarregar")
	} else {
		if comando == "produtos" {
			counts, err := utilsdb.CountProductsByPlatform()
			if err != nil {
				panic(err)
			}
			for plataforma, info := range counts {
				botview.SendReply(bot, update.Message.Chat.ID, fmt.Sprintf("🔒Logins no estoque:\n🔑 Plataforma: [%s]\n💡 Quantidade de logins: [%d]\n💵Preço: [%dR$]", plataforma, info.Quantidade, info.Preco))
			}
		}else if produtosArray(comando, plataformas){
			var preco int
			func ()  {
				if err := db.QueryRow("SELECT preco FROM products WHERE plataforma = ?", comando).Scan(&preco); err != nil{
					log.Fatal(err)
				}
				getSaldo, err := utilsdb.GetUserSaldo(int64(userID))
				if err != nil{
					botview.SendReply(bot, update.Message.Chat.ID, "Não foi possivel comprar, erro no seu saldo atual")
					return
				}else{
					saldoAtual = int(getSaldo)
					if !SaldoSuficiente(saldoAtual, preco){
						botview.SendReply(bot, update.Message.Chat.ID, "Saldo insuficiente 👎!\nUse /saldo para ver seu saldo atual")
						return
					}else{
						login, err := utilsdb.ComprarLogin(comando)
						if err != nil{
							fmt.Println(comando)
							botview.SendReply(bot, update.Message.Chat.ID, fmt.Sprintf("Erro ao comprar login: ", err))
							return
						}
						func ()  {
							_, err := db.Exec("UPDATE users SET saldo = saldo - ? WHERE id = ?", preco, userID)	
							if err != nil{
								log.Fatal(err)
							}
						}()
						botview.SendReply(bot, update.Message.Chat.ID, fmt.Sprintf("═══════ - ════════\nCompra realizada com sucesso!!\n├👤 Email: %s\n└🔑 Senha: %s\nAproveite seu acesso!!\n═══════ - ════════", login.Email, login.Senha))
					}
				}
				

			}()
			return
		}else{
			botview.SendReply(bot, update.Message.Chat.ID, "Este login não esta no estoque ou use /comprar para mais infos")
			return
		}
	}
}

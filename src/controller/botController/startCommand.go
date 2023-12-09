package botcontroller

import (
	"fmt"
	utilsdb "telebotgo/src/db/utilsDB"
	botview "telebotgo/src/view/bot_view"
	"log"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)



func StartCommand(msg *tg.Message, bot *tg.BotAPI, userID int )  {
	
	saldoAtual, err := utilsdb.GetUserSaldo(int64(userID))
	if err != nil{
		log.Fatal(err)
	}
	reply := fmt.Sprintf("Bem vindo(a), %s 💸\n\nVocê está na %s, store. 🪙\n\nSeu saldo atual é de %dR$ 🪪\n\nDigite /help para ver os comandos. 🛠", msg.From.FirstName, bot.Self.UserName, saldoAtual)
	botview.SendReply(bot, msg.Chat.ID, reply)
}
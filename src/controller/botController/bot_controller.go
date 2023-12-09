package botcontroller

import (
	"log"
	utilsdb "telebotgo/src/db/utilsDB"
	botview "telebotgo/src/view/bot_view"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleMessage(update tg.Update, bot *tg.BotAPI, mpToken string) {
	userID := update.Message.From.ID
	utilsdb.InsertNewUser(int64(userID))
	var saldoAtual int64
	msg := update.Message
	if getsaldoAtual, err := utilsdb.GetUserSaldo(int64(msg.From.ID)); err != nil {
		log.Fatal(err)
	} else {
		saldoAtual += getsaldoAtual
	}
	if msg == nil {
		return
	}
	switch msg.Command() {
	case "start":
		StartCommand(msg, bot, userID)
	case "help", "h":
		HelpCommand(bot, msg)
	case "recarregar":
		RecarregarCommand(update, bot, mpToken, userID, saldoAtual)
	case "comprar":
		ComprarCommand(update, bot, userID)
	case "saldo":
		SaldoAtual(update, bot, userID)
	default:
		reply := "Digite /help ou /h para ver os comandos da loja"
		botview.SendReply(bot, msg.Chat.ID, reply)
	}
}

func HandleCallbackQuery(callback *tg.CallbackQuery) {

}

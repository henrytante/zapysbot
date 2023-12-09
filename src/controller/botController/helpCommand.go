package botcontroller

import (
	botview "telebotgo/src/view/bot_view"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)



func HelpCommand(bot *tg.BotAPI, msg *tg.Message)  {
	reply := "🌟 Bem-vindo ao Menu Principal 🌟\n\n/recarregar 🔄 Adicionar saldo\n/comprar 🛒 Ver produtos à venda\n/saldo  💵Ver saldo atual\n"
	botview.SendReply(bot, msg.Chat.ID, reply)
}
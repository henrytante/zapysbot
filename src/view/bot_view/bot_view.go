package botview

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendReply(bot *tg.BotAPI, chatID int64, message string)  {
	msg := tg.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil{
		log.Println("Erro ao enviar menssagem: ", err)
	}
}

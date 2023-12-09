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
func SendReplyMK(bot *tg.BotAPI, chatID int64, message string)  {
	msg := tg.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	if err != nil{
		log.Println("Erro ao enviar menssagem: ", err)
	}
}

func SendInlineKeyboard(bot *tg.BotAPI, chatID int64, message string, inlineButtons [][]tg.InlineKeyboardButton) {
	inlineKeyboardMarkup := tg.NewInlineKeyboardMarkup(inlineButtons...)
	msg := tg.NewMessage(chatID, message)
	msg.ReplyMarkup = inlineKeyboardMarkup

	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Erro ao enviar mensagem: ", err)
	}
}

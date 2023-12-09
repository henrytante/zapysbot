package botcontroller

import (
	"fmt"
	"log"
	utilsdb "telebotgo/src/db/utilsDB"
	botview "telebotgo/src/view/bot_view"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SaldoAtual(update tg.Update, bot *tg.BotAPI, userID int)  {
	getSaldo, err := utilsdb.GetUserSaldo(int64(userID))
	if err != nil{
		log.Fatal(err)
	}
	botview.SendReply(bot, update.Message.Chat.ID, fmt.Sprintf("🪪Seu saldo atual: %d", getSaldo))
}
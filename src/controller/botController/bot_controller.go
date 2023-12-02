package botcontroller

import (
	"fmt"
	"log"
	"strconv"
	"telebotgo/payment"
	"telebotgo/payment/status"
	utilsdb "telebotgo/src/db/utilsDB"
	botview "telebotgo/src/view/bot_view"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleMessage(update tg.Update, bot *tg.BotAPI, mpToken string) {
	userID := update.Message.From.ID
	utilsdb.InsertNewUser(int64(userID))
	var saldoAtual int64
	msg := update.Message
	breve := "esta função sera adicionada em breve, aguarde. Dev: @h1000dev"
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

		reply := fmt.Sprintf("Olá! Seja Bem-Vindo ao %s 🛒\n\nSeu saldo atual é: %vRS\n\n Use /h ou /help", bot.Self.UserName, saldoAtual)
		botview.SendReply(bot, msg.Chat.ID, reply)
	case "help", "h":
		reply := "MENU: \n/recarregar Adicionar saldo\n/comprar Ver produtos a venda\n/contas Ver contas ja compradas"
		botview.SendReply(bot, msg.Chat.ID, reply)
	case "recarregar":
		argumento := msg.CommandArguments()
		if argumento == "" {
			botview.SendReply(bot, msg.Chat.ID, "Digite a quantia que deseja recarregar\nEX:\n/recarregar 10\nOBS: o mínimo é 5RS")
			return
		}

		valor, err := strconv.Atoi(argumento)
		if err != nil {
			botview.SendReply(bot, msg.Chat.ID, "Valor incompatível")
			return
		}
		if valor < 5{
			botview.SendReply(bot, msg.Chat.ID, "O valor minimo é 5RS")
			return
		}
		qrCode, err := payment.PIX(valor, mpToken, userID)
		if err != nil {
			botview.SendReply(bot, msg.Chat.ID, "Erro ao gerar QR_Code")
		} else {
			botview.SendReply(bot, msg.Chat.ID, fmt.Sprintf("Pix gerado com sucesso✅\n\n💠 QR_Code gerado com sucesso: %s", qrCode.QRCode))
			go status.Status(mpToken, userID)
		}

	case "comprar":
		botview.SendReply(bot, msg.Chat.ID, breve)
	case "contas":

		botview.SendReply(bot, msg.Chat.ID, breve)
	default:
		reply := "Digite /help ou /h para ver os comandos da loja"
		botview.SendReply(bot, msg.Chat.ID, reply)
	}
}
func HandleCallbackQuery(callback *tg.CallbackQuery) {

}

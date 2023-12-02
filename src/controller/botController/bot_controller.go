package botcontroller

import (
	"fmt"
	"log"
	"telebotgo/payment"

	utilsdb "telebotgo/src/db/utilsDB"
	botview "telebotgo/src/view/bot_view"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleMessage(update tg.Update, bot *tg.BotAPI, mpToken string) {
	userID := update.Message.From.ID
	utilsdb.InsertNewUser(int64(userID))
	var saldoAtual int64
	msg := update.Message
	breve := "Esta função será adicionada em breve. Aguarde. Dev: @h1000dev"

	if getsaldoAtual, err := utilsdb.GetUserSaldo(int64(msg.From.ID)); err != nil {
		log.Fatal(err)
	} else {
		saldoAtual += getsaldoAtual
	}

	if msg == nil {
		return
	}
	ch := make(chan bool)
	switch msg.Command() {
	case "start":
		reply := fmt.Sprintf("🔷 Olá %s\n\n🔘 Seja bem-vindo, você está no %s🛒\n💸Seu saldo atual é: %vRS\n\nUse /h ou /help", msg.From.FirstName, bot.Self.UserName, saldoAtual)
		botview.SendReply(bot, msg.Chat.ID, reply)

	case "help", "h":
		reply := "MENU: \n/recarregar Adicionar saldo\n/comprar Ver produtos à venda\n/contas Ver contas já compradas"
		botview.SendReply(bot, msg.Chat.ID, reply)

	case "recarregar":
		valor := 0.01
		qrCode, err := payment.PIX(valor, mpToken, userID, ch)
		if err != nil {
			botview.SendReply(bot, msg.Chat.ID, "Erro ao gerar QR_Code")
		} else {
			botview.SendReply(bot, msg.Chat.ID, fmt.Sprintf("Pix gerado com sucesso✅\n💠 QR_Code pix gerado `%s`", qrCode.QRCode))
			go func() {
				paymentApproved := <-ch
				if paymentApproved {
					botview.SendReply(bot, msg.Chat.ID, "Pagamento aprovado")
				} else {
					botview.SendReply(bot, msg.Chat.ID, "Erro no pagamento")
				}
			}()
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
	// Lógica para manipular ações de callback
}

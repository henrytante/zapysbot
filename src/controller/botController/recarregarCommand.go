package botcontroller

import (
	"fmt"
	"log"
	"strconv"
	"telebotgo/payment"
	utilsdb "telebotgo/src/db/utilsDB"
	botview "telebotgo/src/view/bot_view"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RecarregarCommand(update tg.Update, bot *tg.BotAPI, mpToken string, userID int, saldoAtual int64) {
	
	ch := make(chan bool)
	if len(update.Message.CommandArguments()) > 0 {
		valorStr := update.Message.CommandArguments()
		valor, err := strconv.ParseFloat(valorStr, 64)
		if err != nil {
			reply := "Valor invalido. Por favor, insira apenas numeros"
			msg := tg.NewMessage(update.Message.Chat.ID, reply)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		if valor < 4{
			botview.SendReply(bot, update.Message.Chat.ID, "O valor minimo de recarga é 5")
			return
		}
		qrCode, err := payment.PIX(valor, mpToken, userID, ch)
		if err != nil {
			botview.SendReply(bot, update.Message.Chat.ID, "Erro ao gerar QR_Code")
			return
		} else {
			botview.SendReplyMK(bot, update.Message.Chat.ID, fmt.Sprintf("✨ Pix gerado com sucesso! ✨\n\n💠 Pix copia a cola gerado:\n\n`%s`\n\n⏰ Pix válido por 15 minutos ⏰", qrCode.QRCode))
			go func() {
				paymentApproved := <-ch
				if paymentApproved {
					botview.SendReply(bot, update.Message.Chat.ID, fmt.Sprintf("💸 Pagamento aprovado, %vR$ foram adicionados ao seu saldo", valor))
					if err = utilsdb.AddSaldo(int64(userID), int64(valor)); err != nil {
						botview.SendReply(bot, update.Message.Chat.ID, "Não foi possivel adicionar seu saldo, entre em contato com o adm da loja")
						log.Fatal(err)
					}
					func ()  {
						getsaldoAtual, err := utilsdb.GetUserSaldo(int64(userID))
						if err != nil{
							log.Fatal(err)
						}else{
							saldoAtual += getsaldoAtual
						}
					}()
					msg := update.Message
					StartCommand(msg, bot, userID)
				} else {
					botview.SendReply(bot, update.Message.Chat.ID, "⚠️Pix expirado!⚠️")
				}
			}()
		}
	}else{
		reply := "Insira a quantia que deseja recarregar, EX:\n\n/recarregar 5\n\nOBS: O valor minimo é 5R$"
		msg := tg.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.Send(msg)
		if err != nil{
			log.Fatal(err)
		}
	}
}

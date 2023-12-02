package main

import (
	
	"log"
	botcontroller "telebotgo/src/controller/botController"
	"telebotgo/src/db/initdb"
	"telebotgo/utils/clear"
	_ "telebotgo/utils/clear"

	"github.com/fatih/color"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init()  {
	initdb.InitDB()
}

func main()  {
	green := color.New(color.FgGreen).SprintFunc()
	var token string
	var mpToken string
	token = "6314832066:AAF86Zvn1JULp2RiFMg1-u7uzuTZqhgQpoM"
	mpToken = "APP_USR-1823582536914556-120102-b0a99302f5ea4c3fc632de575bfa5d23-1559876250"
	//fmt.Print("\nDigite seu token do telegram: ")
	//fmt.Scan(&token)
	//fmt.Print("\nDigite seu token do mercado pago: ")
	//fmt.Scan(&mpToken)
	bot, err := tg.NewBotAPI(token)
	if err != nil{
		log.Fatal(err)
	}
	clear.Clear()
	log.Println(green("Autenticado com BOT: ", bot.Self.UserName))
	updateConfig := tg.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil{
		log.Fatal(err)
	}
	for update := range updates{
		if update.CallbackQuery != nil{
			botcontroller.HandleCallbackQuery(update.CallbackQuery)
		}else if update.Message != nil{
			botcontroller.HandleMessage(update, bot, mpToken)
		}
	}
}
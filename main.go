package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/paavill/awesome-tagger-bot/balancer"
	bt "github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/cases/process_update"
	con "github.com/paavill/awesome-tagger-bot/repository/mongo"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	con.Init()
	bt.Init()
	process_update.Init()

	log.Printf("Authorized on account %s", bt.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bt.Bot.GetUpdatesChan(u)

	balancer.Run()

	for update := range updates {
		balancer.ReceiveUpdate(update)
	}
}

package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/paavill/awesome-tagger-bot/balancer"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/cases/process_update"
	"github.com/paavill/awesome-tagger-bot/repository/mongo"
	"github.com/paavill/awesome-tagger-bot/scheduler"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mongo.Init()
	bot.Init()
	process_update.Init()

	log.Printf("Authorized on account %s", bot.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.Bot.GetUpdatesChan(u)

	scheduler.Run()
	balancer.Run()

	for update := range updates {
		balancer.ReceiveUpdate(update)
	}
}

package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/paavill/awesome-tagger-bot/balancer"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/context"
	"github.com/paavill/awesome-tagger-bot/domain/cases/process_update"
	domainContext "github.com/paavill/awesome-tagger-bot/domain/context"
	domainLogger "github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/logger"
	"github.com/paavill/awesome-tagger-bot/repository/mongo"
	"github.com/paavill/awesome-tagger-bot/scheduler"
	"github.com/paavill/awesome-tagger-bot/services"
	"github.com/paavill/awesome-tagger-bot/services/kandinsky"
)

func main() {

	time.Local = time.UTC

	kandinsky := kandinsky.NewService("", "", "")
	servicesBuilder := services.NewBuilder().Kandinsky(kandinsky)
	logger := logger.New(domainLogger.Debug)

	ctx, err := context.NewBuilder().
		Logger(logger).
		ServicesBuilder(servicesBuilder).
		Build()

	if err != nil {
		panic(fmt.Sprintf("unable to create context: %s", err))
	}

	domainContext.Set(ctx)

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

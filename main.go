package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/paavill/awesome-tagger-bot/balancer"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/context"
	"github.com/paavill/awesome-tagger-bot/domain/cases/process_update"
	domainLogger "github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/logger"
	"github.com/paavill/awesome-tagger-bot/repository/mongo"
	"github.com/paavill/awesome-tagger-bot/scheduler"
	"github.com/paavill/awesome-tagger-bot/services"
	"github.com/paavill/awesome-tagger-bot/services/kandinsky"
)

func main() {

	time.Local = time.UTC
	logger := logger.New(domainLogger.Debug)

	bot, err := bot.Init(logger)
	if err != nil {
		panic(fmt.Sprintf("unable to init bot: %s", err))
	}

	kandinsky := kandinsky.NewService("", "", "")
	servicesBuilder := services.
		NewBuilder().
		Kandinsky(kandinsky).
		Bot(bot)

	ctx, err := context.NewBuilder().
		Logger(logger).
		ServicesBuilder(servicesBuilder).
		Build()
	if err != nil {
		panic(fmt.Sprintf("unable to create context: %s", err))
	}

	mongo.Init()
	process_update.Init()

	ctx.Logger().Info("Authorized on account %s", ctx.Services().Bot().Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := ctx.Services().Bot().GetUpdatesChan(u)

	scheduler.Run()
	balancer.Run()

	for update := range updates {
		balancer.ReceiveUpdate(update)
	}
}

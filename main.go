package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/paavill/awesome-tagger-bot/balancer"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/config"
	"github.com/paavill/awesome-tagger-bot/context"
	"github.com/paavill/awesome-tagger-bot/domain/cases/command/clear_new_cache"
	"github.com/paavill/awesome-tagger-bot/domain/cases/command/generate_image"
	"github.com/paavill/awesome-tagger-bot/domain/cases/command/news"
	"github.com/paavill/awesome-tagger-bot/domain/cases/command/reset"
	"github.com/paavill/awesome-tagger-bot/domain/cases/command/settings"
	"github.com/paavill/awesome-tagger-bot/domain/cases/process_update"
	domainLogger "github.com/paavill/awesome-tagger-bot/domain/logger"
	domainSm "github.com/paavill/awesome-tagger-bot/domain/state_machine"
	"github.com/paavill/awesome-tagger-bot/logger"
	"github.com/paavill/awesome-tagger-bot/repository/mongo"
	"github.com/paavill/awesome-tagger-bot/scheduler"
	"github.com/paavill/awesome-tagger-bot/services"
	"github.com/paavill/awesome-tagger-bot/services/geonode"
	"github.com/paavill/awesome-tagger-bot/services/kandinsky"
	"github.com/paavill/awesome-tagger-bot/state_machine"
)

func main() {

	time.Local = time.UTC
	logger := logger.New(domainLogger.Debug)

	bot, err := bot.Init(logger)
	if err != nil {
		panic(fmt.Sprintf("unable to init bot: %s", err))
	}

	geonode := geonode.New(config.Env.Geonode.Host)
	kandinsky := kandinsky.NewService(config.Env.Kandinsky.Host, config.Env.Kandinsky.Key, config.Env.Kandinsky.Secret)
	servicesBuilder := services.
		NewBuilder().
		Kandinsky(kandinsky).
		GetProxy(geonode).
		Bot(bot)

	connection := mongo.New()

	ctx, err := context.NewBuilder().
		Logger(logger).
		ServicesBuilder(servicesBuilder).
		Connection(connection).
		Build()
	if err != nil {
		panic(fmt.Sprintf("unable to create context: %s", err))
	}

	ctx.Logger().Info("Authorized on account %s", ctx.Services().Bot().Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := ctx.Services().Bot().GetUpdatesChan(u)

	stateMachine := state_machine.New([]domainSm.State{
		clear_new_cache.New(),
		generate_image.New(),
		news.New(),
		reset.New(),
		settings.New(),
	}, []domainSm.Preprocessor{
		process_update.Run,
	})

	scheduler.Run(ctx)
	balancer.Run(stateMachine)

	for update := range updates {
		balancer.ReceiveUpdate(ctx, update)
	}
}

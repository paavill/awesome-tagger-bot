package services

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Builder interface {
	Kandinsky(kandinsky Kandinsky) Builder
	GetProxy(getProxy GetProxy) Builder
	Bot(bot *tgbotapi.BotAPI) Builder
	Build() (Services, error)
}

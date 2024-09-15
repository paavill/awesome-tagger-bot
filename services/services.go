package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

type svr struct {
	kandinsky services.Kandinsky
	bot       *tgbotapi.BotAPI
}

func (s *svr) Kandinsky() services.Kandinsky {
	return s.kandinsky
}

func (s *svr) Bot() *tgbotapi.BotAPI {
	return s.bot
}

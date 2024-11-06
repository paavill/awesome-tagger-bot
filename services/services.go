package services

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

type svr struct {
	kandinsky services.Kandinsky
	getProxy  services.GetProxy
	bot       *tgbotapi.BotAPI
}

func (s *svr) Kandinsky() services.Kandinsky {
	return s.kandinsky
}

func (s *svr) GetProxy() services.GetProxy {
	return s.getProxy
}

func (s *svr) Bot() *tgbotapi.BotAPI {
	return s.bot
}

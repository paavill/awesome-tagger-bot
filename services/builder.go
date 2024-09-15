package services

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

func NewBuilder() services.Builder {
	return &builder{
		services: &svr{},
	}
}

type builder struct {
	services *svr
}

func (b *builder) Kandinsky(kandinsky services.Kandinsky) services.Builder {
	b.services.kandinsky = kandinsky
	return b
}

func (b *builder) Bot(bot *tgbotapi.BotAPI) services.Builder {
	b.services.bot = bot
	return b
}

func (b *builder) Build() (services.Services, error) {
	if b.services.kandinsky == nil {
		return nil, fmt.Errorf("services: kandinsky is nil")
	}

	if b.services.bot == nil {
		return nil, fmt.Errorf("services: bot is nil")
	}

	return b.services, nil
}

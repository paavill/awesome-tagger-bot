package services

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Services interface {
	Kandinsky() Kandinsky
	Bot() *tgbotapi.BotAPI
}

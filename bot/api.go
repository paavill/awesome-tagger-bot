package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var Bot *tgbotapi.BotAPI

func Init() {
	initBot()
	initCommands()
}

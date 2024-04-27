package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func initCommands() {
	Bot.Send(tgbotapi.SetMyCommandsConfig{
		Commands: []tgbotapi.BotCommand{
			tgbotapi.BotCommand{
				Command:     "/reset",
				Description: "Запрос кнопки \"Поделиться именем\"",
			},
			tgbotapi.BotCommand{
				Command:     "/clear_cash",
				Description: "(Пока не работает) Очистка пользователей из кеша бота",
			},
		},
	})
}

package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func initCommands() {
	_, err := Bot.Send(tgbotapi.SetMyCommandsConfig{
		Commands: []tgbotapi.BotCommand{
			tgbotapi.BotCommand{
				Command:     "/reset",
				Description: "Запрос кнопки \"Поделиться именем\"",
			},
			tgbotapi.BotCommand{
				Command:     "/news",
				Description: "Какой сегодня день?",
			},
			tgbotapi.BotCommand{
				Command:     "/settings",
				Description: "Настройки",
			},
			tgbotapi.BotCommand{
				Command:     "/generate_image",
				Description: "Сгенерировать изображение",
			},
			/*tgbotapi.BotCommand{
				Command:     "/clear_cash",
				Description: "(Пока не работает) Очистка пользователей из кеша бота",
			},*/
		},
	})
	log.Println("Init command err:", err.Error())
}

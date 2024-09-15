package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
)

func initCommands(logger logger.Logger, bot *tgbotapi.BotAPI) {
	_, err := bot.Send(tgbotapi.SetMyCommandsConfig{
		Commands: []tgbotapi.BotCommand{
			{
				Command:     "/reset",
				Description: "Запрос кнопки \"Поделиться именем\"",
			},
			{
				Command:     "/news",
				Description: "Какой сегодня день?",
			},
			{
				Command:     "/settings",
				Description: "Настройки",
			},
			{
				Command:     "/generate_image",
				Description: "Сгенерировать изображение",
			},
			{
				Command:     "/clear_news_cache",
				Description: "Очистить кеш новостей",
			},
			/*tgbotapi.BotCommand{
				Command:     "/clear_cash",
				Description: "(Пока не работает) Очистка пользователей из кеша бота",
			},*/
		},
	})
	logger.Error("init bot command err: %s", err)
}

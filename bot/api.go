package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
)

func Init(logger logger.Logger) (*tgbotapi.BotAPI, error) {
	bot, err := initBot(logger)
	if err != nil {
		return nil, fmt.Errorf("wile init bot due: %s", err)
	}
	initCommands(logger, bot)
	return bot, nil
}

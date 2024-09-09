package bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/config"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
)

func initBot(logger logger.Logger) (*tgbotapi.BotAPI, error) {
	tokenFromFileRaw, err := os.ReadFile(config.Env.Bot.TokenFilename)
	if err != nil {
		logger.Error("reading Token file due: %s", err)
	}
	tokenFromFile := string(tokenFromFileRaw)
	tokenFromFile = strings.ReplaceAll(tokenFromFile, "\n", "")

	tokenFromEnv := config.Env.Bot.Token

	if tokenFromFile != "" && tokenFromEnv != "" && tokenFromFile != tokenFromEnv {
		return nil, fmt.Errorf("token from file and env var are different")
	}

	var token string
	if tokenFromFile != "" {
		token = tokenFromFile
	} else {
		token = tokenFromEnv
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = config.Env.Bot.Debug

	return bot, nil
}

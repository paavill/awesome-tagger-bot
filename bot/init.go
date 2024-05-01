package bot

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/config"
)

func initBot() {
	tokenFromFileRaw, err := os.ReadFile(config.Env.Bot.TokenFilename)
	if err != nil {
		log.Println("Error reading Token file " + err.Error())
	}
	tokenFromFile := string(tokenFromFileRaw)
	tokenFromFile = strings.ReplaceAll(tokenFromFile, "\n", "")

	tokenFromEnv := config.Env.Bot.Token

	if tokenFromFile != "" && tokenFromEnv != "" && tokenFromFile != tokenFromEnv {
		log.Panic("Token from file and env var are different")
	}

	var token string
	if tokenFromFile != "" {
		token = tokenFromFile
	} else {
		token = tokenFromEnv
	}

	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	Bot.Debug = config.Env.Bot.Debug
}

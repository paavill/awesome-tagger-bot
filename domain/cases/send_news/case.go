package send_news

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/cases/get_news"
)

func Run(chatId int64) {
	title, news, err := get_news.Run(chatId)
	if err != nil {
		return
	}
	sendToBot(chatId, title, news)
}

func prepareText(title string, news []string) string {
	result := ""
	result += title + "\n"
	for _, n := range news {
		result += n + "\n"
	}
	return result
}

func sendToBot(chatId int64, title string, news []string) {
	text := prepareText(title, news)
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.Bot.Send(msg)
	if err != nil {
		log.Panicln("Error while sending news " + err.Error())
	}
}

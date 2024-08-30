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

	limitedNews := make([]string, 0, 15)
	for i := 0; i < len(limitedNews) && i < len(news); i++ {
		limitedNews = append(limitedNews, news[i])
	}

	for _, n := range limitedNews {
		result += n + "\n"
	}

	result += "Информация взята с сайта https://kakoysegodnyaprazdnik.ru/"

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

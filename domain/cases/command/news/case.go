package news

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_news"
)

func Run(chatId int64, message *tgbotapi.Message) {
	if message == nil {
		return
	}

	if message.Text != "/news" && message.Text != "/news@"+bot.Bot.Self.UserName {
		return
	}
	send_news.Run(chatId)
}

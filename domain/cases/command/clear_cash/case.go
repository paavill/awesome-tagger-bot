package clear_cash

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/models"
)

func Run(chat *models.Chat, message *tgbotapi.Message) {
	if message == nil {
		return
	}

	if chat != nil && (message.Text == "/clear_cash" || message.Text == "/clear_cash@"+bot.Bot.Self.UserName) {
		//chat.ClearCash = false
	}
}

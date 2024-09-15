package send_greetings

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

func Run(ctx context.Context, tgChatId int64) error {

	chat, err := ctx.Connection().Chat().GetByTgId(tgChatId)
	if err != nil {
		return fmt.Errorf("error while getting chat by tg id: %s", err)
	}

	msg := tgbotapi.NewMessage(tgChatId, `
Привет 😊

Пожалуйста нажми <b>Поделиться именем</b>.

<i>Если ты сделаешь это,</i>
<i>твои друзья смогут</i>
<i>тегать тебя с помощью @all</i>
				`)
	//msg.ReplyToMessageID = update.Message.MessageIDs
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.InlineKeyboardButton{
					Text:         "Поделиться именем",
					CallbackData: &chat.UuidCallback,
				},
			},
		},
	}
	_, err = ctx.Services().Bot().Send(msg)
	if err != nil {
		return fmt.Errorf("error while sending message due: %s", err)
	}
	return nil
}

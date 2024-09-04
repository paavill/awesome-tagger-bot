package generate_image

import (
	"bytes"
	"image"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_images"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

func Run(message *tgbotapi.Message) {
	if message == nil {
		return
	}

	if message.Text == "/generate_image" || message.Text == "/generate_image@"+bot.Bot.Self.UserName {
		raw, err := os.ReadFile("./test.jpg")
		if err != nil {
			context.Get().Logger().Error(err.Error())
			return
		}

		img, _, err := image.Decode(bytes.NewBuffer(raw))
		if err != nil {
			context.Get().Logger().Error(err.Error())
			return
		}

		err = send_images.Run(message.Chat.ID, []*image.Image{&img, &img, &img, &img, &img})
		if err != nil {
			context.Get().Logger().Error(err.Error())
			return
		}
	}
}

package send_images

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

func Run(ctx context.Context, chatId int64, caption string, images []*image.Image) error {
	if len(images) == 0 {
		return fmt.Errorf("no images to send")
	}

	photoConfigs := []any{}
	for i, img := range images {
		writer := bytes.Buffer{}

		if err := jpeg.Encode(&writer, *img, nil); err != nil {
			return fmt.Errorf("error while encoding image due: " + err.Error())
		}

		uuid := uuid.New()
		fileBytes := tgbotapi.FileBytes{
			Name:  uuid.String() + ".jpeg",
			Bytes: writer.Bytes(),
		}

		config := tgbotapi.NewInputMediaPhoto(fileBytes)
		if i == 0 {
			config.Caption = caption
		}
		photoConfigs = append(photoConfigs, config)
	}

	mediaGroup := tgbotapi.NewMediaGroup(chatId, photoConfigs)
	messages, err := ctx.Services().Bot().SendMediaGroup(mediaGroup)
	if err != nil {
		return fmt.Errorf("error while sending media group due: %s", err)
	}
	ctx.Logger().Info("media group send: returned messages: %d", len(messages))
	return nil
}

package send_news

import (
	"fmt"
	"image"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/cases/get_news"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_images"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

func Run(ctx context.Context, chatId int64, autoSend bool) {
	title, news, err := get_news.Run(ctx, chatId, autoSend)
	if err != nil {
		ctx.Logger().Error("[send_news] error while getting news: %s", err)
		return
	}

	if len(news) >= 1 {
		img := get_news.GetImage()
		if img == nil {
			tmpImg, err := ctx.Services().Kandinsky().GenerateImage(news[0])
			if err != nil {
				ctx.Logger().Error("[send_news] error while generating image: %s", err)
			}
			img = tmpImg
			get_news.SetImage(img)
		}

		if img != nil {
			err = send_images.Run(ctx, chatId, "Изображение первого праздника\n создано Kandinsky-им", []*image.Image{img})
			if err != nil {
				ctx.Logger().Error("[send_news] error while sending image: %s", err)
			}
		}
	}

	err = sendToBot(ctx, chatId, title, news)
	if err != nil {
		ctx.Logger().Error("[send_news] error while sending news: %s", err)
	}

}

func prepareText(title string, news []string) string {
	result := ""
	result += title + "\n"

	limitedNews := make([]string, 7)
	for i := 0; i < len(limitedNews) && i < len(news); i++ {
		limitedNews[i] = news[i]
	}

	for _, n := range limitedNews {
		result += n + "\n"
	}

	result += "\nИнформация взята с сайта: https://kakoysegodnyaprazdnik.ru/"

	return result
}

func sendToBot(ctx context.Context, chatId int64, title string, news []string) error {
	text := prepareText(title, news)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := ctx.Services().Bot().Send(msg)
	if err != nil {
		return fmt.Errorf("error while sending news: %s", err)
	}
	return nil
}

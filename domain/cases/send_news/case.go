package send_news

import (
	"fmt"
	"image"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/cases/get_news"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_images"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

func Run(ctx context.Context, chatId int64) {
	title, news, err := get_news.Run(ctx, chatId)
	if err != nil {
		ctx.Logger().Error("[send_news] error while getting news: %s", err)
		return
	}
	if len(news) >= 3 {
		req := strings.Join(news[0:4], "\n")
		req = "Нарисуй картинку, которая будет отражать все это, ГЛАВНОЕ НЕ РИСУЙ ЛЮДЕЙ!!!: \n" + req
		img, err := ctx.Services().Kandinsky().GenerateImage(req)
		if err != nil {
			ctx.Logger().Error("[send_news] error while generating image: %s", err)
		} else {
			err = send_images.Run(ctx, chatId, "", []*image.Image{img})
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

	limitedNews := make([]string, 15)
	for i := 0; i < len(limitedNews) && i < len(news); i++ {
		limitedNews[i] = news[i]
	}

	for _, n := range limitedNews {
		result += n + "\n"
	}

	result += "Информация взята с сайта https://kakoysegodnyaprazdnik.ru/"

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

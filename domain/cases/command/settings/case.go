package settings

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/models"
	"github.com/paavill/awesome-tagger-bot/scheduler"
)

var (
	nul             = "nil"
	root            = "root"
	settingsToday   = "settings_today"
	incrementHour   = "increment_hour"
	incrementMinute = "increment_minute"
	decrementHour   = "decrement_hour"
	decrementMinute = "decrement_minute"
	save            = "save"
	onOffChange     = "on_off_change"
	back            = "back"
	markUps         = map[string]tgbotapi.InlineKeyboardMarkup{
		root: tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         "Какой сегодня день",
						CallbackData: &settingsToday,
					},
				},
			},
		},
		settingsToday: tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         "+",
						CallbackData: &incrementHour,
					},
					tgbotapi.InlineKeyboardButton{
						Text:         "+",
						CallbackData: &incrementMinute,
					},
				},
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         "12",
						CallbackData: &nul,
					},
					tgbotapi.InlineKeyboardButton{
						Text:         "0",
						CallbackData: &nul,
					},
				},
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         "-",
						CallbackData: &decrementHour,
					},
					tgbotapi.InlineKeyboardButton{
						Text:         "-",
						CallbackData: &decrementMinute,
					},
				},
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         "on",
						CallbackData: &onOffChange,
					},
				},
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         "Сохранить",
						CallbackData: &save,
					},
				},
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         "Назад",
						CallbackData: &back,
					},
				},
			},
		},
	}
)

func Run(chatId int64, message *tgbotapi.Message) {
	if message == nil {
		return
	}

	if message.Text != "/settings" && message.Text != "/settings@"+bot.Bot.Self.UserName {
		return
	}

	nmsgc := tgbotapi.NewMessage(chatId, "Настройки\n(только не жми на + и - слишком быстро, я могу устать)")
	nmsgc.ReplyMarkup = markUps[root]

	_, err := bot.Bot.Send(nmsgc)
	if err != nil {
		log.Println("Error sending settings message: ", err)
	}
}

func ProcessCallBack(chatId int64, callbackQuery *tgbotapi.CallbackQuery) {
	if callbackQuery == nil {
		log.Println("Error while processing callback: callbackQuery is nil")
		return
	}

	data := callbackQuery.Data

	message := callbackQuery.Message
	if message == nil {
		log.Println("Error while processing callback: message is nil")
		return
	}

	messageId := message.MessageID

	switch data {
	case incrementHour:
		markup := message.ReplyMarkup
		h, _ := strconv.Atoi(markup.InlineKeyboard[1][0].Text)

		if h+1 < 24 {
			markup.InlineKeyboard[1][0].Text = strconv.Itoa(h + 1)
		} else {
			markup.InlineKeyboard[1][0].Text = strconv.Itoa(0)
		}
		sendMarkupUpdate(chatId, messageId, markup, callbackQuery.ID)
	case incrementMinute:
		markup := message.ReplyMarkup
		m, _ := strconv.Atoi(markup.InlineKeyboard[1][1].Text)

		if m+1 < 60 {
			markup.InlineKeyboard[1][1].Text = strconv.Itoa(m + 1)
		} else {
			markup.InlineKeyboard[1][1].Text = strconv.Itoa(0)
		}
		sendMarkupUpdate(chatId, messageId, markup, callbackQuery.ID)
	case decrementHour:
		markup := message.ReplyMarkup
		h, _ := strconv.Atoi(markup.InlineKeyboard[1][0].Text)

		if h-1 > 0 {
			markup.InlineKeyboard[1][0].Text = strconv.Itoa(h - 1)
		} else {
			markup.InlineKeyboard[1][0].Text = strconv.Itoa(23)
		}
		sendMarkupUpdate(chatId, messageId, markup, callbackQuery.ID)
	case decrementMinute:
		markup := message.ReplyMarkup
		m, _ := strconv.Atoi(markup.InlineKeyboard[1][1].Text)

		if m-1 > 0 {
			markup.InlineKeyboard[1][1].Text = strconv.Itoa(m - 1)
		} else {
			markup.InlineKeyboard[1][1].Text = strconv.Itoa(59)
		}
		sendMarkupUpdate(chatId, messageId, markup, callbackQuery.ID)
	case save:
		markup := message.ReplyMarkup
		m, _ := strconv.Atoi(markup.InlineKeyboard[1][1].Text)
		h, _ := strconv.Atoi(markup.InlineKeyboard[1][0].Text)
		oo := markup.InlineKeyboard[3][0].Text

		schedule := oo == "on"

		oldS, err := scheduler.GetNewsSettingById(chatId)

		if err != nil {
			s := &models.NewsSettings{
				ChatId:   chatId,
				Hour:     h,
				Minute:   m,
				Schedule: schedule,
			}
			scheduler.Process(s)
		} else {
			oldS.Hour = h
			oldS.Minute = m
			oldS.Schedule = schedule

			scheduler.Process(oldS)
		}
		qr := tgbotapi.NewCallback(callbackQuery.ID, "Сохранил!)")
		bot.Bot.Send(qr)
	case back:
		v, _ := markUps[root]
		sendMarkupUpdate(chatId, messageId, &v, callbackQuery.ID)
	case onOffChange:
		markup := message.ReplyMarkup
		m := markup.InlineKeyboard[3][0].Text

		if m == "on" {
			markup.InlineKeyboard[3][0].Text = "off"
		} else {
			markup.InlineKeyboard[3][0].Text = "on"
		}

		sendMarkupUpdate(chatId, messageId, markup, callbackQuery.ID)
	case root:
		markup := markUps[root]
		sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
	case settingsToday:
		markup := markUps[settingsToday]
		oldS, err := scheduler.GetNewsSettingById(chatId)
		if err != nil {
			sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
			return
		}
		markup.InlineKeyboard[1][1].Text = fmt.Sprint(oldS.Minute)
		markup.InlineKeyboard[1][0].Text = fmt.Sprint(oldS.Hour)
		sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
	default:
		log.Println("Unsupported callback data: ", data)
		qr := tgbotapi.NewCallback(callbackQuery.ID, "O.O")
		bot.Bot.Send(qr)
	}
}

func sendMarkupUpdate(chatId int64, messageId int, markup *tgbotapi.InlineKeyboardMarkup, qid string) {
	_, err := bot.Bot.Send(tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   messageId,
			ReplyMarkup: markup,
		},
	})
	if err != nil {
		log.Println("Error sending settings message: ", err)
		alert := tgbotapi.NewCallbackWithAlert(qid, "Меня вот так ругает TG (не надо так быстро): "+err.Error()+" (это в секундах)")
		bot.Bot.Send(alert)
	}
}

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
	nul           = "nil"
	root          = "root"
	hours         = "hours"
	minutes       = "minutes"
	settingsToday = "settings_today"
	selectHour    = "select_hour"
	selectMinute  = "select_minute"
	save          = "save"
	onOffChange   = "on_off_change"
	back          = "back"

	hoursArray   = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}
	minutesArray = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "53", "54", "55", "56", "57", "58", "59"}

	markUps = map[string]tgbotapi.InlineKeyboardMarkup{
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
						Text:         "12",
						CallbackData: &selectHour,
					},
					tgbotapi.InlineKeyboardButton{
						Text:         "0",
						CallbackData: &selectMinute,
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
		hours: tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
			},
		},
		minutes: tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
				[]tgbotapi.InlineKeyboardButton{},
			},
		},
	}
)

func init() {
	for i, v := range hoursArray {
		subArrayIndex := i / 6
		callBackData := "h-" + v
		markUps[hours].InlineKeyboard[subArrayIndex] = append(markUps[hours].InlineKeyboard[subArrayIndex], tgbotapi.InlineKeyboardButton{
			Text:         v,
			CallbackData: &callBackData,
		})
	}
	for i, v := range minutesArray {
		subArrayIndex := i / 6
		callBackData := "m-" + v
		markUps[minutes].InlineKeyboard[subArrayIndex] = append(markUps[minutes].InlineKeyboard[subArrayIndex], tgbotapi.InlineKeyboardButton{
			Text:         v,
			CallbackData: &callBackData,
		})
	}
}

func Run(chatId int64, message *tgbotapi.Message) {
	if message == nil {
		return
	}

	if message.Text != "/settings" && message.Text != "/settings@"+bot.Bot.Self.UserName {
		return
	}

	nmsgc := tgbotapi.NewMessage(chatId, "Настройки")
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
	case selectHour:
		markup := markUps[hours]
		sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
	case selectMinute:
		markup := markUps[minutes]
		sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
	case save:
		markup := message.ReplyMarkup
		m, _ := strconv.Atoi(markup.InlineKeyboard[0][1].Text)
		h, _ := strconv.Atoi(markup.InlineKeyboard[0][0].Text)
		oo := markup.InlineKeyboard[2][0].Text

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
		m := markup.InlineKeyboard[2][0].Text

		if m == "on" {
			markup.InlineKeyboard[2][0].Text = "off"
		} else {
			markup.InlineKeyboard[2][0].Text = "on"
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
		markup.InlineKeyboard[0][1].Text = fmt.Sprint(oldS.Minute)
		markup.InlineKeyboard[0][0].Text = fmt.Sprint(oldS.Hour)
		sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
	default:
		for _, v := range hoursArray {
			if "h-"+v == data {
				markup := markUps[settingsToday]
				markup.InlineKeyboard[0][0].Text = v
				sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
				return
			}
		}

		for _, v := range minutesArray {
			if "m-"+v == data {
				markup := markUps[settingsToday]
				markup.InlineKeyboard[0][1].Text = v
				sendMarkupUpdate(chatId, messageId, &markup, callbackQuery.ID)
				return
			}
		}

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

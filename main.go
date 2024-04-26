package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/paavill/awesome-tagger-bot/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chat struct {
	ID        int64
	Users     map[string]struct{}
	New       bool
	ClearCash bool
}

var chats = map[int64]*Chat{}
var collection *mongo.Collection
var bot *tgbotapi.BotAPI
var uid string = uuid.NewString()

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	b, err := os.ReadFile(config.Env.Bot.TokenFilename)
	if err != nil {
		log.Panic(err)
	}
	token := string(b)
	token = strings.ReplaceAll(token, "\n", "")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = config.Env.Bot.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		processUpdate(update)
		if strings.Contains(update.Message.Text, "@all") {
			tags := []string{}
			for u, _ := range chats[update.Message.Chat.ID].Users {
				tags = append(tags, "@"+u)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				`
Простите, пошумлю:
`+strings.Join(tags, "\n"))
			bot.Send(msg)
		}
	}
}

func processUpdate(update tgbotapi.Update) {
	cbq := update.CallbackQuery
	chat := update.FromChat()
	id := chat.ID
	initBotIfNeed(id)
	clearCashCommand(id, "")
	resetCommand(id, "")
	processChat(id)
	callbackProcess(cbq, id)
}

func initBotIfNeed(id int64) {
	if _, ok := chats[id]; !ok {
		chats[id] = &Chat{
			ID:        id,
			Users:     map[string]struct{}{},
			New:       true,
			ClearCash: false,
		}
		log.Printf("Chat with ID %d added", id)
	}
}

func processChat(id int64) {
	chat, ok := chats[id]
	if !ok {
		panic("Не должно такого происходить")
	}

	if chat.ClearCash {
		chat.Users = map[string]struct{}{}
		chat.ClearCash = false
	}

	if chat.New {
		initUsers(id)
		chat.New = false
	}
}

func initUsers(id int64) {
	if ch, ok := chats[id]; ok {
		msg := tgbotapi.NewMessage(ch.ID, `
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
						CallbackData: &uid,
					},
				},
			},
		}
		bot.Send(msg)
	}
}

func callbackProcess(q *tgbotapi.CallbackQuery, chatId int64) {
	data := q.Data
	username := q.From.UserName
	if data == uid {
		chats[chatId].Users[username] = struct{}{}
	}
}

func clearCashCommand(id int64, command string) {
	if ch, ok := chats[id]; ok && command == "/clear_cash" {
		ch.ClearCash = true
	}
}

func resetCommand(id int64, command string) {
	if ch, ok := chats[id]; ok && command == "/reset" {
		ch.New = true
	}
}

package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chat struct {
	ID    int64
	Users []string
}

var chats = map[int64]*Chat{}
var collection *mongo.Collection

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://:@localhost:27017"))
	if err != nil {
		log.Panic(err)
	}

	collection = client.Database("tagger").Collection("users")

	t, ok := os.LookupEnv("BOT_TOKEN_FILENAME")
	if !ok {
		log.Panic("BOT_TOKEN_FILENAME is empty")
	}
	b, err := os.ReadFile(t)
	if err != nil {
		log.Panic(err)
	}
	token := string(b)
	token = strings.ReplaceAll(token, "\n", "")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("Username [%s]", update.Message.From.UserName)

			if _, ok := chats[update.Message.Chat.ID]; !ok {
				chats[update.Message.Chat.ID] = &Chat{
					ID:    update.Message.Chat.ID,
					Users: []string{},
				}
				log.Printf("Added chat %d", update.Message.Chat.ID)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
Привет 😊

Пожалуйста нажми <b>Поделиться именем</b>.

<i>Если ты сделаешь это,</i>
<i>твои друзья смогут</i>
<i>тегать тебя с помощью @all</i>
				`)
				//msg.ReplyToMessageID = update.Message.MessageID
				msg.ParseMode = "HTML"

				msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
					OneTimeKeyboard: true,
					Keyboard: [][]tgbotapi.KeyboardButton{
						{
							{
								Text: "Поделиться именем",
							},
						},
					},
				}
				bot.Send(msg)
			}
		}

		if strings.Contains(update.Message.Text, "@all") {
			tags := []string{}
			for _, u := range chats[update.Message.Chat.ID].Users {
				tags = append(tags, "@"+u)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				`
Простите, пошумлю:
`+strings.Join(tags, "\n"))
			bot.Send(msg)
		}

		if update.Message.Text == "Поделиться именем" {
			chats[update.Message.Chat.ID].Users = append(chats[update.Message.Chat.ID].Users, update.Message.From.UserName)
		}
	}
}

package process_update

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	bt "github.com/paavill/awesome-tagger-bot/bot"
	"github.com/paavill/awesome-tagger-bot/domain/models"
	"github.com/paavill/awesome-tagger-bot/repository/mongo"
)

var (
	chats   = map[int64]*models.Chat{}
	ownName = "awesome_tagger_bot"
)

// TODO загружать подругому
func Init() {
	chs, err := mongo.Chats().FindAll()
	if err != nil {
		log.Panic("This shouldn't happen")
	}
	for _, ch := range chs {
		chats[ch.Id] = &ch
	}
}

func Run(update tgbotapi.Update) {
	chat := update.FromChat()
	if chat == nil {
		chat = &update.MyChatMember.Chat
	}
	id := chat.ID
	chatName := chat.Title
	cbq := update.CallbackQuery

	initChatIfNeed(id, chatName)
	clearCashCommand(id, "")
	if update.Message != nil {
		resetCommand(id, update.Message.Text)
	}
	processChat(id)
	callbackProcess(cbq, id)
	processTagAll(update)

	if _, ok := chats[id]; !ok {
		log.Panic("This shouldn't happen")
	}

	ch := chats[id]
	if ch.MongoId == "" {
		chat, err := mongo.Chats().Insert(*ch)
		if err != nil {
			log.Printf("Error while inserting chat %d to mongo", ch.Id)
		}
		chats[chat.Id] = &chat
	} else {
		err := mongo.Chats().Update(*ch)
		if err != nil {
			log.Printf("Error while updating chat %d with mongo", ch.Id)
		}
	}
}

func initChatIfNeed(id int64, chatName string) {
	if _, ok := chats[id]; !ok {
		chats[id] = &models.Chat{
			Id:           id,
			ChatName:     chatName,
			Users:        map[string]struct{}{},
			New:          true,
			ClearCash:    false,
			UuidCallback: uuid.NewString(),
		}
		log.Printf("Chat with ID %d added", id)
	} else if chats[id].ChatName == "" {
		chats[id].ChatName = chatName
	}
}

func processChat(id int64) {
	chat, ok := chats[id]
	if !ok {
		panic("This shouldn't happen")
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
		msg := tgbotapi.NewMessage(ch.Id, `
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
						CallbackData: &ch.UuidCallback,
					},
				},
			},
		}
		bt.Bot.Send(msg)
	}
}

func callbackProcess(q *tgbotapi.CallbackQuery, chatId int64) {
	if q == nil {
		return
	}
	data := q.Data
	user := q.From
	if user == nil {
		return
	}
	username := user.UserName
	if data == chats[chatId].UuidCallback {
		chats[chatId].Users[username] = struct{}{}
		log.Printf("User %s shared name in chat %d", username, chatId)
		callBackConfig := tgbotapi.NewCallbackWithAlert(q.ID, "Спасибо, теперь я тебя знаю☺")
		bt.Bot.Send(callBackConfig)
	}
}

func processTagAll(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if strings.Contains(update.Message.Text, "@all") {
		tags := []string{}
		for u, _ := range chats[update.Message.Chat.ID].Users {
			tags = append(tags, "@"+u)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			`
Простите, пошумлю:
`+strings.Join(tags, "\n"))
		bt.Bot.Send(msg)
	}
}

func clearCashCommand(id int64, command string) {
	if ch, ok := chats[id]; ok && command == "/clear_cash" {
		ch.ClearCash = false
	}
}

func resetCommand(id int64, command string) {
	if ch, ok := chats[id]; ok && command == "/reset@awesome_tagger_bot" {
		ch.New = true
	}
}

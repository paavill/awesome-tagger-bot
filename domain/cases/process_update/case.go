package process_update

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_greetings"
	"github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/models"
)

func Run(ctx context.Context, update *tgbotapi.Update) error {
	chat := update.FromChat()
	if chat == nil {
		chat = &update.MyChatMember.Chat
	}
	id := chat.ID
	cbq := update.CallbackQuery

	err := initChatIfNeed(ctx, update)
	if err != nil {
		return fmt.Errorf("error while init chat: %s", err)
	}

	err = processShareUsername(ctx, cbq, id)
	if err != nil {
		return fmt.Errorf("error while processing share username: %s", err)
	}

	err = processTagAll(ctx, id, update)
	if err != nil {
		return fmt.Errorf("error while processing tag all: %s", err)
	}

	return nil
}

func initChatIfNeed(ctx context.Context, update *tgbotapi.Update) error {
	tgChat := update.FromChat()
	if tgChat == nil {
		return fmt.Errorf("chat is nil")
	}
	chat, err := ctx.Connection().Chat().GetByTgId(tgChat.ID)
	if err == nil {
		cq := update.CallbackQuery
		if cq == nil {
			return nil
		}

		usr := cq.From
		if usr == nil {
			return nil
		}

		if cq.Data == chat.UuidCallback {
			chat.Users[usr.UserName] = struct{}{}
			err = ctx.Connection().Chat().Update(chat)
			if err != nil {
				return fmt.Errorf("error while updating chat: %s", err)
			}
		}

		return nil
	}
	chat = &models.Chat{
		Id:           tgChat.ID,
		ChatName:     tgChat.Title,
		UuidCallback: uuid.New().String(),
		Users:        map[string]struct{}{},
	}
	err = ctx.Connection().Chat().Insert(chat)
	if err != nil {
		return fmt.Errorf("error while inserting chat to mongo: %s", err)
	}

	err = send_greetings.Run(ctx, tgChat.ID)
	if err != nil {
		return fmt.Errorf("error while sending greetings: %s", err)
	}

	return nil
}

func processShareUsername(ctx context.Context, q *tgbotapi.CallbackQuery, chatId int64) error {
	if q == nil {
		return nil //fmt.Errorf("[process_update] callback query is nil")
	}
	data := q.Data
	user := q.From
	if user == nil {
		return fmt.Errorf("[process_update] user is nil")
	}
	username := user.UserName
	chat, err := ctx.Connection().Chat().GetByTgId(chatId)
	if err != nil {
		return fmt.Errorf("error while getting chat %d from mongo: %s", chatId, err)
	}
	if data == chat.UuidCallback {
		chat.Users[username] = struct{}{}
		ctx.Logger().Info("User %s shared name in chat %d", username, chatId)
		callBackConfig := tgbotapi.NewCallback(q.ID, "Спасибо, теперь я тебя знаю☺")
		ctx.Services().Bot().Send(callBackConfig)
	}
	return nil
}

func processTagAll(ctx context.Context, chatId int64, update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil //fmt.Errorf("[process_update] message is nil")
	}
	if strings.Contains(update.Message.Text, "@all") {
		tags := []string{}
		chat, err := ctx.Connection().Chat().GetByTgId(chatId)
		if err != nil {
			return fmt.Errorf("error while getting chat %d from mongo: %s", chatId, err)
		}
		for u, _ := range chat.Users {
			if u == update.Message.From.UserName {
				continue
			}
			tags = append(tags, "@"+u)
		}
		msg := tgbotapi.NewMessage(chatId,
			`
Простите, пошумлю:
`+strings.Join(tags, "\n"))
		allMsg, err := ctx.Services().Bot().Send(msg)
		if err != nil {
			return fmt.Errorf("[process_update] error while sending message: %s", err)
		}

		time.Sleep(1 * time.Second)

		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    update.Message.Chat.ID,
				MessageID: allMsg.MessageID,
			},
			Text: "Я пошумел, всех вызвал!\nИ прибрал за собой😅",
		}
		_, err = ctx.Services().Bot().Send(edit)
		if err != nil {
			return fmt.Errorf("[process_update] error while editing message: %s", err)
		}
	}
	return nil
}

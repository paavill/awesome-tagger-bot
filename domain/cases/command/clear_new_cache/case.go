package clear_new_cache

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/cases/get_news"
	"github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/state_machine"
)

func New() state_machine.State {
	return &state{}
}

type state struct {
	state_machine.Dumper
}

func (s *state) ProcessCallbackRequest(ctx context.Context, callback *tgbotapi.CallbackQuery) (state_machine.ProcessResponse, error) {
	return nil, nil
}

func (s *state) ProcessMessage(ctx context.Context, message *tgbotapi.Message) (state_machine.ProcessResponse, error) {
	err := Run(ctx, message)
	return nil, err
}

func Run(ctx context.Context, message *tgbotapi.Message) error {
	if message == nil {
		return fmt.Errorf("[clear_new_cache] message is nil")
	}

	selfName := ctx.Services().Bot().Self.UserName
	if message.Text != "/clear_news_cache" && message.Text != "/clear_news_cache@"+selfName {
		return nil
	}

	ctx.Logger().Info("[clear_new_cache] start")
	defer ctx.Logger().Info("[clear_new_cache] end")

	get_news.ClearCache()

	chat := message.Chat

	callBackConfig := tgbotapi.NewMessage(chat.ID, "Очистил👍")
	ctx.Services().Bot().Send(callBackConfig)
	return nil
}

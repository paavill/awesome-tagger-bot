package clear_cash

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	ctx.Logger().Info("[clear_cash] start")
	defer ctx.Logger().Info("[clear_cash] end")
	if message == nil {
		return fmt.Errorf("[clear_cash] message is nil")
	}

	selfName := ctx.Services().Bot().Self.UserName
	if message.Text != "/clear_cash" && message.Text != "/clear_cash@"+selfName {
		return nil
	}

	messageChat := message.Chat
	if messageChat == nil {
		return fmt.Errorf("[clear_cash] chat is nil")
	}

	chat, err := ctx.Connection().Chat().GetByTgId(messageChat.ID)
	if err != nil {
		return fmt.Errorf("[clear_cash] error get chat due: %s", err)
	}

	chat.Users = map[string]struct{}{}

	err = ctx.Connection().Chat().Update(chat)
	if err != nil {
		return fmt.Errorf("[clear_cash] error update chat due: %s", err)
	}
	return nil
}

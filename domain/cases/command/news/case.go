package news

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_news"
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
		return fmt.Errorf("[news] message is nil")
	}

	selfName := ctx.Services().Bot().Self.UserName
	if message.Text != "/news" && message.Text != "/news@"+selfName {
		return nil
	}

	ctx.Logger().Info("[news] start")
	defer ctx.Logger().Info("[news] end")

	messageChat := message.Chat
	if messageChat == nil {
		return fmt.Errorf("[news] message chat is nil")
	}

	send_news.Run(ctx, messageChat.ID)

	return nil
}

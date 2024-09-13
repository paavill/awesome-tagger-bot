package reset

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/cases/send_greetings"
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
	ctx.Logger().Info("[reset] start")
	defer ctx.Logger().Info("[reset] end")
	if message == nil {
		return fmt.Errorf("[reset] message is nil")
	}

	selfName := ctx.Services().Bot().Self.UserName
	if message.Text != "/reset" || message.Text != "/reset@"+selfName {
		return nil
	}

	messageChat := message.Chat
	if messageChat == nil {
		return fmt.Errorf("[reset] message chat is nil")
	}

	err := send_greetings.Run(ctx, messageChat.ID)
	if err != nil {
		return fmt.Errorf("[reset] error while sending greetings due: %s", err)
	}
	return nil
}

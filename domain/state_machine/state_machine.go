package state_machine

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

type StateMachine interface {
	GetInitStates() []State
	Process(context.Context, tgbotapi.Update) error
}

type ProcessResponse interface {
	States() []State
	NeedAddInitStates() bool
}

type State interface {
	ProcessCallbackRequest(context.Context, *tgbotapi.CallbackQuery) (ProcessResponse, error)
	ProcessMessage(context.Context, *tgbotapi.Message) (ProcessResponse, error)
	Dump() string
}

package state_machine

import (
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

type Dumper struct {
}

func (s *Dumper) Dump() string {
	stack := []byte{}
	runtime.Stack(stack, true)
	return string(stack)
}

type StateMachine interface {
	Process(context.Context, tgbotapi.Update) error
}

func NewProcessResponse(states ...State) ProcessResponse {
	return &processResponse{states: states}
}

type processResponse struct {
	states []State
}

func (r *processResponse) States() []State {
	return r.states
}

type ProcessResponse interface {
	States() []State
}

type State interface {
	ProcessCallbackRequest(context.Context, *tgbotapi.CallbackQuery) (ProcessResponse, error)
	ProcessMessage(context.Context, *tgbotapi.Message) (ProcessResponse, error)
	Dump() string
}

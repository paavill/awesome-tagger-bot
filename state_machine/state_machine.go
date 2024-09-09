package state_machine

import (
	"errors"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

type processResponse interface {
	Chattable() tgbotapi.Chattable
	MediaGroup() *tgbotapi.MediaGroupConfig
	States() []state
}

type state interface {
	processCallbackRequest(context.Context, *tgbotapi.CallbackQuery) (processResponse, error)
	processMessage(context.Context, *tgbotapi.Message) (processResponse, error)
	dump() string
}

type stateMachine struct {
	mux           *sync.Mutex
	initStates    []state
	currentStates map[int64][]state
}

func (sm *stateMachine) Process(ctx context.Context, update tgbotapi.Update) error {
	var chatId *int64 = nil

	message := update.Message
	callback := update.CallbackQuery

	if message != nil && callback != nil {
		return fmt.Errorf("message and callback are not nil at the same time")
	}

	if message != nil {
		chatId = &message.Chat.ID
	}

	if callback != nil {
		chatId = &callback.Message.Chat.ID
	}

	sm.mux.Lock()
	states, ok := sm.currentStates[*chatId]
	if !ok {
		states = sm.initStates
		sm.currentStates[*chatId] = states
	}
	sm.mux.Unlock()

	responses, errorStates, err := processStates(ctx, update, states)
	if err != nil || len(errorStates) > 0 {
		sm.mux.Lock()
		defer sm.mux.Unlock()

		msg := tgbotapi.NewMessage(*chatId, "Уууупс, что-то пошло не так...")
		ctx.Services().Bot().Send(msg)

		if len(errorStates) > 0 {
			ctx.Logger().Critical("there are error states")
			dump := ""
			for _, state := range errorStates {
				dump += state.dump() + "\n"
			}
			_, err := ctx.Services().Bot().Send(tgbotapi.NewMessage(*chatId, dump))
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
		}

		return err
	}

	newStates := []state{}
	for _, response := range responses {
		if response.Chattable() != nil {
			_, err := ctx.Services().Bot().Send(response.Chattable())
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
		}
		if response.MediaGroup() != nil {
			_, err := ctx.Services().Bot().SendMediaGroup(*response.MediaGroup())
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
		}
		newStates = append(newStates, response.States()...)
	}

	sm.mux.Lock()
	sm.currentStates[*chatId] = newStates
	sm.mux.Unlock()

	return nil
}

func processStates(ctx context.Context, update tgbotapi.Update, states []state) ([]processResponse, []state, error) {
	resultResponses := []processResponse{}
	errorStates := []state{}

	for _, currentState := range states {
		localResponses := []processResponse{}
		localErrors := []error{}

		response, err := currentState.processCallbackRequest(ctx, update.CallbackQuery)
		if err != nil {
			localErrors = append(localErrors, fmt.Errorf("error while processing callback request: %s", err))
		} else if response != nil {
			localResponses = append(localResponses, response)
		}

		response, err = currentState.processMessage(ctx, update.Message)
		if err != nil {
			localErrors = append(localErrors, fmt.Errorf("error while processing message: %s", err))
		} else if response != nil {
			localResponses = append(localResponses, response)
		}

		if len(localErrors) != 0 {
			return nil, nil, errors.Join(localErrors...)
		}

		if len(localResponses) != 1 {
			errorStates = append(errorStates, currentState)
		} else {
			resultResponses = append(resultResponses, localResponses[0])
		}
	}

	return resultResponses, errorStates, nil
}

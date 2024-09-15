package state_machine

import (
	"errors"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/state_machine"
)

func New(initStates []state_machine.State, preprocessors []state_machine.Preprocessor) state_machine.StateMachine {
	return &stateMachine{
		mux:           &sync.Mutex{},
		preprocessors: preprocessors,
		initStates:    initStates,
		currentStates: map[int64][]state_machine.State{},
	}
}

type stateMachine struct {
	mux           *sync.Mutex
	preprocessors []state_machine.Preprocessor
	initStates    []state_machine.State
	currentStates map[int64][]state_machine.State
}

func (sm *stateMachine) GetInitStates() []state_machine.State {
	result := []state_machine.State{}
	result = append(result, sm.initStates...)
	return result
}

func (sm *stateMachine) Process(ctx context.Context, update tgbotapi.Update) error {
	for _, p := range sm.preprocessors {
		err := p(ctx, &update)
		if err != nil {
			return err
		}
	}

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
				dump += state.Dump() + "\n"
			}
			_, err := ctx.Services().Bot().Send(tgbotapi.NewMessage(*chatId, dump))
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
			ctx.Logger().Error(dump)
		}

		sm.currentStates[*chatId] = sm.initStates

		return err
	}

	newStates := []state_machine.State{}
	addInit := false
	for _, response := range responses {
		if response == nil {
			continue
		}
		addInit = response.NeedAddInitStates() || addInit
		newStates = append(newStates, response.States()...)
	}

	if addInit {
		newStates = append(newStates, sm.initStates...)
	}

	if len(newStates) != 0 {
		sm.mux.Lock()
		sm.currentStates[*chatId] = newStates
		sm.mux.Unlock()
	} else {
		sm.mux.Lock()
		sm.currentStates[*chatId] = sm.initStates
		sm.mux.Unlock()
	}

	return nil
}

func processStates(ctx context.Context, update tgbotapi.Update, states []state_machine.State) ([]state_machine.ProcessResponse, []state_machine.State, error) {
	resultResponses := []state_machine.ProcessResponse{}
	errorStates := []state_machine.State{}

	for _, currentState := range states {
		localResponses := []state_machine.ProcessResponse{}
		localErrors := []error{}

		if update.CallbackQuery != nil {
			response, err := currentState.ProcessCallbackRequest(ctx, update.CallbackQuery)
			if err != nil {
				localErrors = append(localErrors, fmt.Errorf("error while processing callback request: %s", err))
			} else if response != nil {
				localResponses = append(localResponses, response)
			}
		}

		if update.Message != nil {
			response, err := currentState.ProcessMessage(ctx, update.Message)
			if err != nil {
				localErrors = append(localErrors, fmt.Errorf("error while processing message: %s", err))
			} else if response != nil {
				localResponses = append(localResponses, response)
			}
		}

		if len(localErrors) != 0 {
			return nil, nil, errors.Join(localErrors...)
		}

		if len(localResponses) > 1 {
			errorStates = append(errorStates, currentState)
		} else if len(localResponses) == 1 {
			resultResponses = append(resultResponses, localResponses[0])
		}
	}

	return resultResponses, errorStates, nil
}

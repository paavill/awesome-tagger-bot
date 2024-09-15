package state_machine

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/context"
)

type Preprocessor func(context.Context, *tgbotapi.Update) error

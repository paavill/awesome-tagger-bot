package context

import (
	"github.com/paavill/awesome-tagger-bot/domain/connection"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

type Context interface {
	Services() services.Services
	Connection() connection.Connection
	Logger() logger.Logger
}

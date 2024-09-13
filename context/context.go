package context

import (
	"github.com/paavill/awesome-tagger-bot/domain/connection"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

type ctx struct {
	connection connection.Connection
	services   services.Services
	logger     logger.Logger
}

func (c *ctx) Connection() connection.Connection {
	return c.connection
}

func (c *ctx) Services() services.Services {
	return c.services
}

func (c *ctx) Logger() logger.Logger {
	return c.logger
}

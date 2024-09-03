package context

import (
	"github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

type ctx struct {
	services services.Services
	logger   logger.Logger
}

func (c *ctx) Services() services.Services {
	return c.services
}

func (c *ctx) Logger() logger.Logger {
	return c.logger
}

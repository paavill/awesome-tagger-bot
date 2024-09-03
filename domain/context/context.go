package context

import (
	"github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

var (
	context Context
)

type Context interface {
	Services() services.Services
	Logger() logger.Logger
}

func Set(ctx Context) {
	if context != nil {
		// TODO может быть паника во время выполнения
		ctx.Logger().Error("warning: context already set")
		return
	}
	context = ctx
}

func Get() Context {
	if context == nil {
		// TODO может быть паника во время выполнения
		panic("context not set")
	}
	return context
}

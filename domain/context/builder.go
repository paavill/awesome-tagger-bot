package context

import (
	"github.com/paavill/awesome-tagger-bot/domain/connection"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

type Builder interface {
	ServicesBuilder(services.Builder) Builder
	Logger(logger logger.Logger) Builder
	Connection(connection connection.Connection) Builder
	Build() (Context, error)
}

package context

import (
	"fmt"

	"github.com/paavill/awesome-tagger-bot/domain/connection"
	dc "github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

func NewBuilder() dc.Builder {
	return &builder{}
}

type builder struct {
	services   services.Builder
	logger     logger.Logger
	connection connection.Connection
}

func (b *builder) ServicesBuilder(builder services.Builder) dc.Builder {
	b.services = builder
	return b
}

func (b *builder) Logger(logger logger.Logger) dc.Builder {
	b.logger = logger
	return b
}

func (b *builder) Connection(connection connection.Connection) dc.Builder {
	b.connection = connection
	return b
}

func (b *builder) Build() (dc.Context, error) {
	services, err := b.services.Build()
	if err != nil {
		return nil, err
	}

	if b.logger == nil {
		return nil, fmt.Errorf("logger not set")
	}

	if b.connection == nil {
		return nil, fmt.Errorf("connection not set")
	}

	return &ctx{
		services:   services,
		logger:     b.logger,
		connection: b.connection,
	}, nil
}

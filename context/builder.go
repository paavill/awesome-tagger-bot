package context

import (
	"fmt"

	"github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/logger"
	"github.com/paavill/awesome-tagger-bot/domain/services"
)

func NewBuilder() context.Builder {
	return &builder{}
}

type builder struct {
	services services.Builder
	logger   logger.Logger
}

func (b *builder) ServicesBuilder(builder services.Builder) context.Builder {
	b.services = builder
	return b
}

func (b *builder) Logger(logger logger.Logger) context.Builder {
	b.logger = logger
	return b
}

func (b *builder) Build() (context.Context, error) {
	services, err := b.services.Build()
	if err != nil {
		return nil, err
	}

	if b.logger == nil {
		return nil, fmt.Errorf("logger not set")
	}

	return &ctx{
		services: services,
		logger:   b.logger,
	}, nil
}

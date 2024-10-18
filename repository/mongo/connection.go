package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/paavill/awesome-tagger-bot/config"
	"github.com/paavill/awesome-tagger-bot/domain/connection"
	"github.com/paavill/awesome-tagger-bot/domain/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type conn struct {
	chat         *chatRepo
	newsSettings *newsSettingsRepo
}

func New() connection.Connection {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", config.Env.Mongo.User, config.Env.Mongo.Pass, config.Env.Mongo.Host))) //"mongodb://:@localhost:27017"
	if err != nil {
		log.Panic(err)
	}

	err = c.Ping(ctx, options.Client().ReadPreference)
	if err != nil {
		log.Panic(err)
	}

	return &conn{
		chat:         &chatRepo{makeChatRepo(c)},
		newsSettings: &newsSettingsRepo{makeNewsSettings(c)},
	}
}

func (c *conn) Chat() repositories.Chat {
	return c.chat
}

func (c *conn) NewsSettings() repositories.NewsSettings {
	return c.newsSettings
}

package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/paavill/awesome-tagger-bot/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", config.Env.Mongo.User, config.Env.Mongo.Pass, config.Env.Mongo.Host))) //"mongodb://:@localhost:27017"
	if err != nil {
		log.Panic(err)
	}
	client = c
}

func Chats() ChatsRepo {
	return &chatRepo{
		collection: makeChatRepo(client),
	}
}

func NewsSettings() NewsSettingsRepo {
	return &newsSettingsRepo{
		collection: makeNewsSettings(client),
	}
}

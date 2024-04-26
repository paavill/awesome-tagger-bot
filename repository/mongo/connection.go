package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/paavill/awesome-tagger-bot/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New() (client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", config.Env.Mongo.User, config.Env.Mongo.Pass, config.Env.Mongo.Host))) //"mongodb://:@localhost:27017"
	if err != nil {
		//log.Panic(err)
	}
	return
	//collection = client.Database("tagger").Collection("users")
}

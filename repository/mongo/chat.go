package mongo

import (
	"context"
	"log"

	"github.com/paavill/awesome-tagger-bot/config"
	"github.com/paavill/awesome-tagger-bot/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoChat struct {
	MongoId      primitive.ObjectID `bson:"_id"`
	TgId         int64              `bson:"tg_id"`
	ChatName     string             `bson:"chat_name"`
	UuidCallback string             `bson:"uuid_callback"`
	Users        []string           `bson:"users"`
}

func (c *mongoChat) fromModel(model models.Chat) {
	oid, err := primitive.ObjectIDFromHex(model.MongoId)
	if err != nil || model.MongoId == "" {
		log.Panic("This shouldn't happen")
	}
	c.MongoId = oid
	c.TgId = model.Id
	c.UuidCallback = model.UuidCallback
	c.ChatName = model.ChatName
	users := make([]string, 0, len(model.Users))
	for k := range model.Users {
		users = append(users, k)
	}
	c.Users = users
}

func (c *mongoChat) toModel() models.Chat {
	model := models.Chat{}
	model.Id = c.TgId
	model.MongoId = c.MongoId.Hex()
	model.Users = map[string]struct{}{}
	model.UuidCallback = c.UuidCallback
	model.ChatName = c.ChatName
	for _, u := range c.Users {
		model.Users[u] = struct{}{}
	}
	return model
}

type chatRepo struct {
	collection *mongo.Collection
}

func makeChatRepo(client *mongo.Client) *mongo.Collection {
	collection := client.Database(config.Env.Mongo.Db).Collection("chats")
	return collection
}

func (r *chatRepo) FindAll() ([]*models.Chat, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var chats []mongoChat
	if err = cursor.All(context.TODO(), &chats); err != nil {
		return nil, err
	}

	result := make([]*models.Chat, 0, len(chats))
	for _, mch := range chats {
		model := mch.toModel()
		result = append(result, &model)
	}
	return result, nil
}

func (r *chatRepo) GetById(id string) (*models.Chat, error) {
	res := r.collection.FindOne(context.Background(), &bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}
	mch := &mongoChat{}
	err := res.Decode(mch)
	if err != nil {
		return nil, err
	}
	model := mch.toModel()
	return &model, nil
}

func (r *chatRepo) GetByTgId(id int64) (*models.Chat, error) {
	res := r.collection.FindOne(context.Background(), &bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}
	mch := &mongoChat{}
	err := res.Decode(mch)
	if err != nil {
		return nil, err
	}
	model := mch.toModel()
	return &model, nil
}

func (r *chatRepo) Insert(model *models.Chat) error {
	mch := mongoChat{}
	model.MongoId = primitive.NewObjectID().Hex()
	mch.fromModel(*model)
	_, err := r.collection.InsertOne(context.Background(), mch)
	return err
}

func (r *chatRepo) Update(model *models.Chat) error {
	if model.MongoId == "" {
		log.Panic("This shouldn't happen")
	}
	mch := mongoChat{}
	mch.fromModel(*model)
	_, err := r.collection.ReplaceOne(context.Background(), bson.M{"_id": mch.MongoId}, mch)
	return err
}

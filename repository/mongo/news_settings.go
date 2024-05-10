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

type mongoNewsSettings struct {
	MongoId  primitive.ObjectID `bson:"_id"`
	ChatId   int64              `bson:"chat_id"`
	Hour     int                `bson:"hour"`
	Minute   int                `bson:"minute"`
	Schedule bool               `bson:"schedule"`
}

func (s *mongoNewsSettings) fromModel(model models.NewsSettings) {
	oid, err := primitive.ObjectIDFromHex(model.MongoId)
	if err != nil || model.MongoId == "" {
		log.Panic("This shouldn't happen")
	}
	s.MongoId = oid
	s.ChatId = model.ChatId
	s.Hour = model.Hour
	s.Minute = model.Minute
	s.Schedule = model.Schedule
}

func (s *mongoNewsSettings) toModel() models.NewsSettings {
	model := models.NewsSettings{}
	model.ChatId = s.ChatId
	model.MongoId = s.MongoId.Hex()
	model.Hour = s.Hour
	model.Minute = s.Minute
	model.Schedule = s.Schedule
	return model
}

type NewsSettingsRepo interface {
	Insert(*models.NewsSettings) error
	Update(*models.NewsSettings) error
	FindAll() ([]models.NewsSettings, error)
}

type newsSettingsRepo struct {
	collection *mongo.Collection
}

func makeNewsSettings(client *mongo.Client) *mongo.Collection {
	collection := client.Database(config.Env.Mongo.Db).Collection("news_settings")
	return collection
}

func (r *newsSettingsRepo) FindAll() ([]models.NewsSettings, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var settings []mongoNewsSettings
	if err = cursor.All(context.TODO(), &settings); err != nil {
		return nil, err
	}

	result := make([]models.NewsSettings, 0, len(settings))
	for _, setting := range settings {
		result = append(result, setting.toModel())
	}
	return result, nil
}

func (r *newsSettingsRepo) Insert(model *models.NewsSettings) error {
	mch := mongoNewsSettings{}
	model.MongoId = primitive.NewObjectID().Hex()
	mch.fromModel(*model)
	_, err := r.collection.InsertOne(context.Background(), mch)
	return err
}

func (r *newsSettingsRepo) Update(model *models.NewsSettings) error {
	if model.MongoId == "" {
		log.Panic("This shouldn't happen")
	}
	mch := mongoNewsSettings{}
	mch.fromModel(*model)
	_, err := r.collection.ReplaceOne(context.Background(), bson.M{"_id": mch.MongoId}, mch)
	return err
}

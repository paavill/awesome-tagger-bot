package config

import (
	"log"
	"os"
	"strconv"
)

var Env Config = Config{}

func init() {
	v := getEnv("MONGODB_URI")
	Env.Mongo.Host = v

	v = getEnv("MONGODB_USER")
	Env.Mongo.User = v

	v = getEnv("MONGODB_PASS")
	Env.Mongo.Pass = v

	v = getEnv("MONGODB_DB")
	Env.Mongo.Db = v

	v = getEnv("BOT_TOKEN_FILENAME")
	Env.Bot.TokenFilename = v

	v = getEnv("BOT_TOKEN")
	Env.Bot.Token = v

	v = getEnv("BOT_DEBUG")
	vb, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatal(err)
	}
	Env.Bot.Debug = vb

	v = getEnv("KANDINSKY_HOST")
	Env.Kandinsky.Host = v

	v = getEnv("KANDINSKY_KEY")
	Env.Kandinsky.Key = v

	v = getEnv("KANDINSKY_SECRET")
	Env.Kandinsky.Secret = v
}

func getEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal(key + " is not set")
	}
	return v
}

type Config struct {
	Mongo     MongoConfig
	Bot       BotConfig
	Kandinsky KandinskyConfig
}

type KandinskyConfig struct {
	Host   string
	Key    string
	Secret string
}

type MongoConfig struct {
	Host string
	User string
	Pass string
	Db   string
}

type BotConfig struct {
	TokenFilename string
	Token         string
	Debug         bool
}

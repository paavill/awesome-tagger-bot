package config

import (
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
		//log.Fatal(err)
	}
	Env.Bot.Debug = vb
}

func getEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		//log.Fatal(key + " is not set")
	}
	return v
}

type Config struct {
	Mongo MongoConfig
	Bot   BotConfig
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

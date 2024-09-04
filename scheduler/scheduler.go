package scheduler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paavill/awesome-tagger-bot/domain/cases/send_news"
	"github.com/paavill/awesome-tagger-bot/domain/models"
	"github.com/paavill/awesome-tagger-bot/repository/mongo"
)

var (
	// TODO race condition
	newsQueue = map[int64]struct {
		cancel  context.CancelFunc
		setting *models.NewsSettings
	}{}
)

func Run() {
	settings, err := mongo.NewsSettings().FindAll()
	if err != nil {
		log.Println("Error while Init scheduler: ", err)
	}

	for _, setting := range settings {
		Process(&setting)
	}

	sysCh := make(chan os.Signal)
	signal.Notify(sysCh, syscall.SIGTERM)
	go func() {
		<-sysCh
		stop()
	}()
}

func Process(setting *models.NewsSettings) {
	add(setting)

	var err error
	if setting.MongoId == "" {
		err = mongo.NewsSettings().Insert(setting)
	} else {
		err = mongo.NewsSettings().Update(setting)
	}
	if err != nil {
		log.Println("Error inserting/updating setting to DB: ", err)
	}
}

func GetNewsSettingById(chatId int64) (*models.NewsSettings, error) {
	setting, ok := newsQueue[chatId]
	if !ok {
		return nil, fmt.Errorf("news setting not found for chatId: %d", chatId)
	}
	return setting.setting, nil
}

func stop() {
	for _, v := range newsQueue {
		v.cancel()
	}
}

func add(setting *models.NewsSettings) {

	v, ok := newsQueue[setting.ChatId]
	if ok {
		v.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	newsQueue[setting.ChatId] = struct {
		cancel  context.CancelFunc
		setting *models.NewsSettings
	}{
		cancel:  cancel,
		setting: setting,
	}

	if setting.Schedule {
		go run(ctx, setting)
	}

}

func remove(setting *models.NewsSettings) {
	v, ok := newsQueue[setting.ChatId]
	if ok {
		v.cancel()
	}
}

func run(ctx context.Context, setting *models.NewsSettings) {
	now := time.Now()
	sleepTime := calcSleepTime(setting.Hour, setting.Minute, now.Hour(), now.Minute())
	log.Println(sleepTime.String())
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(sleepTime):
			now := time.Now()
			sleepTime = calcSleepTime(setting.Hour, setting.Minute, now.Hour(), now.Minute())
			send_news.Run(setting.ChatId)
			log.Println("Sending news at", now, "for chat", setting.ChatId)
		}
	}
}

func calcSleepTime(settingHour, settingMinute, nowHour, nowMinute int) time.Duration {

	nowTime := time.Duration(nowHour)*time.Hour + time.Duration(nowMinute)*time.Minute
	settingTime := time.Duration(settingHour)*time.Hour + time.Duration(settingMinute)*time.Minute

	delta := nowTime - settingTime

	if delta < 0 {
		delta *= -1
	} else {
		delta = time.Duration(24)*time.Hour + time.Duration(nowMinute) - delta
	}
	return delta
}

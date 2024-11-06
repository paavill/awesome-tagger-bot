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
	dc "github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/models"
)

var (
	// TODO race condition
	newsQueue = map[int64]struct {
		cancel  context.CancelFunc
		setting *models.NewsSettings
	}{}
)

func Run(ctx dc.Context) {
	settings, err := ctx.Connection().NewsSettings().FindAll()
	if err != nil {
		log.Println("Error while Init scheduler: ", err)
	}

	for _, setting := range settings {
		Process(ctx, setting)
	}

	sysCh := make(chan os.Signal, 1)
	signal.Notify(sysCh, syscall.SIGTERM)
	go func() {
		<-sysCh
		stop(ctx)
	}()
}

func Process(ctx dc.Context, setting *models.NewsSettings) {
	add(ctx, setting)

	var err error
	if setting.MongoId == "" {
		err = ctx.Connection().NewsSettings().Insert(setting)
	} else {
		err = ctx.Connection().NewsSettings().Update(setting)
	}
	if err != nil {
		ctx.Logger().Error("error inserting/updating setting to DB: ", err)
	}
}

func GetNewsSettingById(chatId int64) (*models.NewsSettings, error) {
	setting, ok := newsQueue[chatId]
	if !ok {
		return nil, fmt.Errorf("news setting not found for chatId: %d", chatId)
	}
	return setting.setting, nil
}

func stop(ctx dc.Context) {
	ctx.Logger().Info("stopping scheduler")
	for _, v := range newsQueue {
		v.cancel()
	}
}

func add(dc dc.Context, setting *models.NewsSettings) {

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
		go run(ctx, dc, setting)
	}

}

func run(ctx context.Context, dc dc.Context, setting *models.NewsSettings) {
	now := time.Now()
	sleepTime := calcSleepTime(setting.Hour, setting.Minute, now.Hour(), now.Minute())
	dc.Logger().Info("send news to chat [%d] after [%s]", setting.ChatId, sleepTime.String())
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(sleepTime):
			now := time.Now()
			sleepTime = calcSleepTime(setting.Hour, setting.Minute, now.Hour(), now.Minute())
			send_news.Run(dc, setting.ChatId, true)
			dc.Logger().Info("sending news at", now, "for chat", setting.ChatId)
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

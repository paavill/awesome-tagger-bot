package balancer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paavill/awesome-tagger-bot/domain/cases/process_update"
)

var (
	queue map[int64]chan tgbotapi.Update = map[int64]chan tgbotapi.Update{}
	sysCh chan os.Signal                 = make(chan os.Signal)
	ctx   context.Context                = context.Background()
)

func Run() {
	signal.Notify(sysCh, syscall.SIGTERM)
	c, f := context.WithCancel(ctx)
	ctx = c
	go func(cancel context.CancelFunc) {
		<-sysCh
		cancel()
	}(f)
}

func ReceiveUpdate(update tgbotapi.Update) {
	chat := update.FromChat()
	if chat == nil {
		chat = &update.MyChatMember.Chat
	}
	id := chat.ID
	if ch, ok := queue[id]; ok {
		ch <- update
	} else {
		queue[id] = make(chan tgbotapi.Update)
		go runChanProcessor(id)
		queue[id] <- update
	}
}

func runChanProcessor(id int64) {
	flag := true
	for flag {
		select {
		case <-ctx.Done():
			flag = false
		case u := <-queue[id]:
			process_update.Run(u)
		}
	}
}

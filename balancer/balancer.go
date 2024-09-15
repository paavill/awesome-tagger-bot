package balancer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cd "github.com/paavill/awesome-tagger-bot/domain/context"
	"github.com/paavill/awesome-tagger-bot/domain/state_machine"
)

var (
	queue        map[int64]chan tgbotapi.Update = map[int64]chan tgbotapi.Update{}
	sysCh        chan os.Signal                 = make(chan os.Signal)
	ctx          context.Context                = context.Background()
	stateMachine state_machine.StateMachine
)

func Run(sm state_machine.StateMachine) {
	stateMachine = sm
	signal.Notify(sysCh, syscall.SIGTERM)
	c, f := context.WithCancel(ctx)
	ctx = c
	go func(cancel context.CancelFunc) {
		<-sysCh
		cancel()
	}(f)
}

func ReceiveUpdate(ctx cd.Context, update tgbotapi.Update) {
	chat := update.FromChat()
	if chat == nil {
		chat = &update.MyChatMember.Chat
	}
	id := chat.ID
	if ch, ok := queue[id]; ok {
		ch <- update
	} else {
		queue[id] = make(chan tgbotapi.Update, 1000)
		go runChanProcessor(ctx, id)
		queue[id] <- update
	}
}

func runChanProcessor(domainContext cd.Context, id int64) {
	flag := true
	for flag {
		select {
		case <-ctx.Done():
			flag = false
		case u := <-queue[id]:
			err := stateMachine.Process(domainContext, u)
			if err != nil {
				domainContext.Logger().Error(err.Error())
			}
		}
	}
}

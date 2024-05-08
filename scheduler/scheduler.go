package scheduler

import (
	"context"
	"time"

	"github.com/paavill/awesome-tagger-bot/domain/models"
)

var (
	newsQueue = map[int64]context.Context{}
)

func Add(setting models.NewsSettings) {
	now := time.Now()
	
}

func Remove(setting models.NewsSettings) {

}

func run(ctx context.Context, firstSleep time.Duration) {
	
}

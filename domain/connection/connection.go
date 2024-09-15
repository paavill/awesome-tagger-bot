package connection

import "github.com/paavill/awesome-tagger-bot/domain/repositories"

type Connection interface {
	Chat() repositories.Chat
	NewsSettings() repositories.NewsSettings
}

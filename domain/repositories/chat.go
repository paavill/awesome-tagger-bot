package repositories

import "github.com/paavill/awesome-tagger-bot/domain/models"

type Chat interface {
	GetById(id string) (*models.Chat, error)
	GetByTgId(id int64) (*models.Chat, error)
	Insert(chat *models.Chat) error
	Update(chat *models.Chat) error
	FindAll() ([]*models.Chat, error)
}

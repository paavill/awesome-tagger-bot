package repositories

import "github.com/paavill/awesome-tagger-bot/domain/models"

type NewsSettings interface {
	Insert(*models.NewsSettings) error
	Update(*models.NewsSettings) error
	FindAll() ([]*models.NewsSettings, error)
}

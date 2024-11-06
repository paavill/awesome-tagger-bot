package services

import "github.com/paavill/awesome-tagger-bot/domain/models"

type GetProxy interface {
	GetProxyList() ([]*models.Proxy, error)
}

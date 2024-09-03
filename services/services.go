package services

import "github.com/paavill/awesome-tagger-bot/domain/services"

type svr struct {
	kandinsky services.Kandinsky
}

func (s *svr) Kandinsky() services.Kandinsky {
	return s.kandinsky
}

package models

import "fmt"

type NewsSettings struct {
	ChatId  int64
	MongoId string
	Hour    int
	Minute  int
}

func (n *NewsSettings) Validate() error {
	if n.Hour < 0 || n.Hour > 23 {
		return fmt.Errorf("Час должен быть между 0 и 23")
	}
	if n.Minute < 0 || n.Minute > 59 {
		return fmt.Errorf("Минута должна быть между 0 и 59")
	}
	return nil
}

package models

type Chat struct {
	Id           int64
	ChatName     string
	MongoId      string
	UuidCallback string
	Users        map[string]struct{}
	New          bool
	ClearCash    bool
}

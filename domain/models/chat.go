package models

type Chat struct {
	Id           int64
	MongoId      string
	UuidCallback string
	Users        map[string]struct{}
	New          bool
	ClearCash    bool
}

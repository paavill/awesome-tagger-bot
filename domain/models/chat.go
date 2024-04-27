package models

type Chat struct {
	Id        int64
	MongoId   string
	Users     map[string]struct{}
	New       bool
	ClearCash bool
}
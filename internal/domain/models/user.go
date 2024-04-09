package models

import "time"

type User struct {
	ID                 int64
	Username           string
	Slug               string
	HashedPass         []byte
	CreatedAt          time.Time
	UpdatedAt          time.Time
	LastPasswordChange time.Time
}

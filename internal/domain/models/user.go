package models

import "time"

type User struct {
	ID                 int32
	Email              string
	Username           string
	Slug               string
	HashedPass         []byte
	CreatedAt          time.Time
	UpdatedAt          time.Time
	LastPasswordChange time.Time
}

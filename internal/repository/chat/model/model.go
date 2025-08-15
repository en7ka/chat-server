package model

import (
	"time"
)

type Chat struct {
	ID        int64
	Name      string
	Type      string
	CreatedAt time.Time
	IsDeleted bool
}

type ChatMember struct {
	ID       int64
	ChatID   int64
	UserId   int64
	JoinedAt time.Time
	IsAdmin  bool
}

type Message struct {
	ID         int64
	ChatID     int64
	FromUserID int64
	Text       string
	Timestamp  time.Time
}

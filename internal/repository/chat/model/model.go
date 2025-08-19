package model

import (
	"time"
)

type Chat struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	IsDeleted bool      `db:"is_deleted"`
}

type ChatMember struct {
	ID       int64     `db:"id"`
	ChatID   int64     `db:"chat_id"`
	UserID   int64     `db:"user_id"`
	JoinedAt time.Time `db:"joined_at"`
	IsAdmin  bool      `db:"is_admin"`
}

type Message struct {
	ID         int64     `db:"id"`
	ChatID     int64     `db:"chat_id"`
	FromUserID int64     `db:"from_user_id"`
	Text       string    `db:"text"`
	Timestamp  time.Time `db:"timestamp"`
}

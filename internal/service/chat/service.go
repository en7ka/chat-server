package chat

import (
	repoif "github.com/en7ka/chat-server/internal/repository/repointerface"
)

type serv struct {
	chatRepository repoif.ChatRepository
}

func NewService(chatRepository repoif.ChatRepository) *serv {
	return &serv{chatRepository: chatRepository}
}

package chat

import (
	"github.com/en7ka/chat-server/internal/client/db"
	repoif "github.com/en7ka/chat-server/internal/repository/repointerface"
)

type serv struct {
	chatRepository repoif.ChatRepository
	txManager      db.TxManager
}

func NewService(chatRepository repoif.ChatRepository, txManager db.TxManager) *serv {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}

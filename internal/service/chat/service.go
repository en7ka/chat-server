package chat

import (
	"github.com/en7ka/chat-server/internal/client/db"
	repoinf "github.com/en7ka/chat-server/internal/repository/repointerface"
	servinf "github.com/en7ka/chat-server/internal/service/servinterface"
)

type chatService struct {
	chatRepository repoinf.ChatRepository
	txManager      db.TxManager
}

func NewService(chatRepository repoinf.ChatRepository, txManager db.TxManager) servinf.ChatService {
	return &chatService{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}

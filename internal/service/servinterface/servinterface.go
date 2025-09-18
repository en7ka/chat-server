package servinterface

import (
	"context"

	"github.com/en7ka/chat-server/internal/models"
)

type ChatService interface {
	CreateChat(ctx context.Context, chat *models.Chat) (*models.Chat, error)
	AddMemberToChat(ctx context.Context, member *models.ChatMember) (int64, error)
	SendMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error)
	GetChatById(ctx context.Context, chatId int64) (*models.Chat, error)
}

type Access interface {
	Access(ctx context.Context, path string) error
}

package servinterface

import (
	"context"
	"github.com/en7ka/chat-server/internal/models"
)

type ChatService interface {
	CreateChat(ctx context.Context, chat *models.Chat) (*models.Chat, error)
	AddMemberToChat(ctx context.Context, member *models.ChatMember) (int64, error)
	SendMessage(ctx context.Context, message *models.Message) (int64, error)
	GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error)
}

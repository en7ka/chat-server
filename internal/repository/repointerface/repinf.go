package repointerface

import (
	"context"
	"github.com/en7ka/chat-server/internal/models"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *models.Chat) (int64, error)
	AddMemberToChat(ctx context.Context, member *models.ChatMember) (int64, error)
	SendMessage(ctx context.Context, message *models.Message) (int64, error)
	GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error)
	GetChatById(ctx context.Context, chatId int64) (*models.Chat, error)
}

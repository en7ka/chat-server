package repointerface

import (
	"context"
	"github.com/en7ka/chat-server/internal/repository/chat/model"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *model.Chat) (int64, error)
	AddMemberToChat(ctx context.Context, member *model.ChatMember) (int64, error)
	SendMessage(ctx context.Context, message *model.Message) (int64, error)
	GetChatMessages(ctx context.Context, chatId int64) ([]*model.Message, error)
}

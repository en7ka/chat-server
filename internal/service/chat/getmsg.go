package chat

import (
	"context"
	"errors"

	"github.com/en7ka/chat-server/internal/models"
)

func (c *chatService) GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error) {
	if chatId <= 0 {
		return nil, errors.New("invalid chat ID")
	}

	return c.chatRepository.GetChatMessages(ctx, chatId)
}

package chat

import (
	"context"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (c *chatService) GetChatById(ctx context.Context, chatId int64) (*models.Chat, error) {
	if chatId <= 0 {
		return nil, fmt.Errorf("invalid chat ID: %d", chatId)
	}

	chat, err := c.chatRepository.GetChatById(ctx, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat by id from repository: %w", err)
	}

	return chat, nil
}

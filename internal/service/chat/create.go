package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (c *chatService) CreateChat(ctx context.Context, chat *models.Chat) (*models.Chat, error) {
	if chat == nil {
		return nil, errors.New("chat is nil")
	}

	// 1. Вызываем репозиторий и получаем ID созданного чата
	chatID, err := c.chatRepository.CreateChat(ctx, chat)
	if err != nil {
		// 2. Если репозиторий вернул ошибку, оборачиваем ее и возвращаем
		return nil, fmt.Errorf("failed to create chat in repository: %w", err)
	}

	chat.ID = chatID

	return chat, nil
}

package chat

import (
	"context"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (s *chatService) CreateChat(ctx context.Context, chat *models.Chat) (*models.Chat, error) {
	var createdChat *models.Chat

	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		chatID, err := s.chatRepository.CreateChat(ctx, chat)
		if err != nil {
			return fmt.Errorf("failed to create chat in repository: %w", err)
		}

		createdChat = &models.Chat{
			ID:   chatID,
			Name: chat.Name,
			Type: chat.Type,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdChat, nil
}

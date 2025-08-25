package chat

import (
	"context"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (s *chatService) GetChatById(ctx context.Context, chatId int64) (*models.Chat, error) {
	var chat *models.Chat

	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var err error
		chat, err = s.chatRepository.GetChatById(ctx, chatId)
		if err != nil {
			return fmt.Errorf("failed to get chat by id from repository: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return chat, nil
}

package chat

import (
	"context"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (s *chatService) GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error) {
	var messages []*models.Message

	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var err error
		messages, err = s.chatRepository.GetChatMessages(ctx, chatId)
		if err != nil {
			return fmt.Errorf("failed to get chat messages: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return messages, nil
}

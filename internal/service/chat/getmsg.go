package chat

import (
	"context"
	"errors"
	"github.com/en7ka/chat-server/internal/models"
)

func (s *serv) GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error) {
	if chatId <= 0 {
		return nil, errors.New("invalid chat ID")
	}

	return s.chatRepository.GetChatMessages(ctx, chatId)
}

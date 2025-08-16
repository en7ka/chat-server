package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/en7ka/chat-server/internal/models"
)

func (s *serv) SendMessage(ctx context.Context, msg *models.Message) (*models.Message, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}

	messageId, err := s.chatRepository.SendMessage(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	msg.ID = messageId

	return msg, nil
}

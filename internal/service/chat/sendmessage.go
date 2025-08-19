package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (c *chatService) SendMessage(ctx context.Context, msg *models.Message) (*models.Message, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}

	messageId, err := c.chatRepository.SendMessage(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	msg.ID = messageId

	return msg, nil
}

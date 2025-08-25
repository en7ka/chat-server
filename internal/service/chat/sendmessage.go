package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (c *chatService) SendMessage(ctx context.Context, msgToCreate *models.Message) (*models.Message, error) {
	if msgToCreate == nil {
		return nil, errors.New("message is nil")
	}

	var createdMsg *models.Message

	err := c.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		messageId, txErr := c.chatRepository.SendMessage(ctx, msgToCreate)
		if txErr != nil {
			return fmt.Errorf("failed to send message in repository: %w", txErr)
		}

		createdMsg = &models.Message{
			ID:         messageId,
			ChatID:     msgToCreate.ChatID,
			FromUserID: msgToCreate.FromUserID,
			Text:       msgToCreate.Text,
			Timestamp:  msgToCreate.Timestamp,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdMsg, nil
}

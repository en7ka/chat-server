package chat

import (
	"context"
	"errors"
	"github.com/en7ka/chat-server/internal/converter"
	"github.com/en7ka/chat-server/internal/models"
)

func (s *serv) SendMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	if message == nil {
		return nil, errors.New("message is nil")
	}
	repoMsg := converter.ToRepoMessageFromDomain(message)

	messageID, err := s.chatRepository.SendMessage(ctx, repoMsg)
	if err != nil {
		return nil, err
	}

	message.ID = messageID

	return message, nil
}

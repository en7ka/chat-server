package chat

import (
	"context"
	"errors"
	"github.com/en7ka/chat-server/internal/models"

	"github.com/en7ka/chat-server/internal/converter"
)

func (s *serv) GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error) {
	if chatId <= 0 {
		return nil, errors.New("invalid chat ID")
	}

	repoMessages, err := s.chatRepository.GetChatMessages(ctx, chatId)
	if err != nil {
		return nil, err
	}

	messages := make([]*models.Message, 0, len(repoMessages))
	for _, repoMsg := range repoMessages {
		messages = append(messages, converter.ToDomainMessageFromRepo(repoMsg))
	}

	return messages, nil
}

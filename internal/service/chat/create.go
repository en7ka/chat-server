package chat

import (
	"context"
	"errors"
	"github.com/en7ka/chat-server/internal/converter"
	"github.com/en7ka/chat-server/internal/models"
)

func (s *serv) CreateChat(ctx context.Context, chat *models.Chat) (*models.Chat, error) {
	if chat == nil {
		return nil, errors.New("chat is nil")
	}

	repoChat := converter.ToRepoChatFromDomain(chat)

	chatId, err := s.chatRepository.CreateChat(ctx, repoChat)
	if err != nil {
		return nil, err
	}

	chat.ID = chatId

	return chat, nil
}

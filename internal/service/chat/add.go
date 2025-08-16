package chat

import (
	"context"
	"errors"
	"github.com/en7ka/chat-server/internal/models"
)

func (s *serv) AddMemberToChat(ctx context.Context, member *models.ChatMember) (int64, error) {
	if member == nil {
		return 0, errors.New("member is nil")
	}

	return s.chatRepository.AddMemberToChat(ctx, member)
}

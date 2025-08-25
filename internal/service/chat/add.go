package chat

import (
	"context"
	"fmt"

	"github.com/en7ka/chat-server/internal/models"
)

func (s *chatService) AddMemberToChat(ctx context.Context, member *models.ChatMember) (int64, error) {
	var memberID int64

	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var err error
		memberID, err = s.chatRepository.AddMemberToChat(ctx, member)
		if err != nil {
			return fmt.Errorf("failed to add member in repository: %w", err)
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return memberID, nil
}

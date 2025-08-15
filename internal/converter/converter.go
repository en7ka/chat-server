package converter

import (
	"github.com/en7ka/chat-server/internal/models"
	repoModel "github.com/en7ka/chat-server/internal/repository/chat/model"
)

func ToRepoMessageFromDomain(domainMsg *models.Message) *repoModel.Message {
	return &repoModel.Message{
		ID:         domainMsg.ID,
		ChatID:     domainMsg.ChatID,
		FromUserID: domainMsg.FromUserID,
		Text:       domainMsg.Text,
		Timestamp:  domainMsg.Timestamp,
	}
}

func ToDomainMessageFromRepo(repoMsg *repoModel.Message) *models.Message {
	return &models.Message{
		ID:         repoMsg.ID,
		ChatID:     repoMsg.ChatID,
		FromUserID: repoMsg.FromUserID,
		Text:       repoMsg.Text,
		Timestamp:  repoMsg.Timestamp,
	}
}

func ToDomainChatFromRepo(repoChat *repoModel.Chat) *models.Chat {
	return &models.Chat{
		ID:        repoChat.ID,
		Name:      repoChat.Name,
		Type:      repoChat.Type,
		CreatedAt: repoChat.CreatedAt,
		IsDeleted: repoChat.IsDeleted,
	}
}

func ToRepoChatFromDomain(domainChat *models.Chat) *repoModel.Chat {
	return &repoModel.Chat{
		ID:        domainChat.ID,
		Name:      domainChat.Name,
		Type:      domainChat.Type,
		CreatedAt: domainChat.CreatedAt,
		IsDeleted: domainChat.IsDeleted,
	}
}

func ToRepoChatMemberFromDomain(domainChatMember *models.ChatMember) *repoModel.ChatMember {
	return &repoModel.ChatMember{
		ID:       domainChatMember.ID,
		ChatID:   domainChatMember.ChatID,
		UserId:   domainChatMember.UserId,
		JoinedAt: domainChatMember.JoinedAt,
		IsAdmin:  domainChatMember.IsAdmin,
	}
}

package converter

import (
	"github.com/en7ka/chat-server/internal/models"
	repoModel "github.com/en7ka/chat-server/internal/repository/chat/model"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
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
		UserID:   domainChatMember.UserID,
		JoinedAt: domainChatMember.JoinedAt,
		IsAdmin:  domainChatMember.IsAdmin,
	}
}

func FromProtoCreateChatRequest(req *desc.CreateRequest) *models.Chat {

	chatName := strings.Join(req.GetUsernames(), ", ")
	return &models.Chat{
		Name: chatName,
		Type: models.ChatTypeGroup,
	}
}

func ToProtoChat(domainChat *models.Chat) *desc.Chat {
	return &desc.Chat{
		Id:        domainChat.ID,
		Name:      domainChat.Name,
		CreatedAt: timestamppb.New(domainChat.CreatedAt),
	}
}

func ToProtoMessage(domainMsg *models.Message) *desc.Message {
	return &desc.Message{
		Id:         domainMsg.ID,
		FromUserId: domainMsg.FromUserID,
		Text:       domainMsg.Text,
		Timestamp:  timestamppb.New(domainMsg.Timestamp),
	}
}

func FromProtoAddMemberRequest(req *desc.AddMemberToChatRequest) *models.ChatMember {
	return &models.ChatMember{
		ChatID: req.GetChatId(),
		UserID: req.GetUserId(),
	}
}

func FromProtoSendMessageRequest(req *desc.SendMessageRequest) *models.Message {
	return &models.Message{
		ChatID:     req.GetChatId(),
		FromUserID: req.GetFromUserId(),
		Text:       req.GetText(),
	}
}

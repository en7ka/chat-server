package chat

import (
	"context"

	"github.com/en7ka/chat-server/internal/converter"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) GetChatMessages(ctx context.Context, req *desc.GetMessagesRequest) (*desc.GetMessagesResponse, error) {
	chatId := req.GetChatId()
	if chatId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Chat ID must be a positive number")
	}

	domainMessages, err := c.chatService.GetChatMessages(ctx, chatId)
	if err != nil {
		return nil, err
	}

	protoMessages := make([]*desc.Message, 0, len(domainMessages))
	for _, domainMsg := range domainMessages {
		protoMessages = append(protoMessages, converter.ToProtoMessage(domainMsg))
	}

	// 4. Формирование ответа
	response := &desc.GetMessagesResponse{
		Messages: protoMessages,
	}

	return response, nil
}

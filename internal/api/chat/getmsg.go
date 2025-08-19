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

	chatMessages, err := c.chatService.GetChatMessages(ctx, chatId)
	if err != nil {
		return nil, err
	}

	protoMessages := make([]*desc.Message, 0, len(chatMessages))
	for _, chatMsg := range chatMessages {
		protoMessages = append(protoMessages, converter.ToProtoMessage(chatMsg))
	}

	// 4. Формирование ответа
	response := &desc.GetMessagesResponse{
		Messages: protoMessages,
	}

	return response, nil
}

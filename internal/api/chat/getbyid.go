package chat

import (
	"context"

	"github.com/en7ka/chat-server/internal/converter"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) GetChat(ctx context.Context, req *desc.GetChatRequest) (*desc.GetChatResponse, error) {
	chatId := req.GetId()
	if chatId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Chat ID must be a positive number")
	}

	domainChat, err := c.chatService.GetChatById(ctx, chatId)
	if err != nil {
		return nil, err
	}

	protoChat := converter.ToProtoChat(domainChat)

	response := &desc.GetChatResponse{
		Chat: protoChat,
	}

	return response, nil
}

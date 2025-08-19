package chat

import (
	"context"
	"strings"

	"github.com/en7ka/chat-server/internal/converter"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*desc.SendMessageResponse, error) {
	if req.GetChatId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id must be greater than zero")
	}

	if req.GetFromUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "from user id must be greater than zero")
	}

	if strings.TrimSpace(req.GetText()) == "" {
		return nil, status.Error(codes.InvalidArgument, "message text cannot be empty")
	}

	messageId := converter.FromProtoSendMessageRequest(req)
	createMsg, err := c.chatService.SendMessage(ctx, messageId)
	if err != nil {
		return nil, err
	}

	return &desc.SendMessageResponse{
		Message: converter.ToProtoMessage(createMsg),
	}, nil
}

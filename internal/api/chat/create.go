package chat

import (
	"context"
	"log"

	"github.com/en7ka/chat-server/internal/converter"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Controller) CreateChat(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if len(req.GetUsernames()) == 0 {
		log.Printf("CreateChat: No usernames provided")
		return nil, status.Error(codes.InvalidArgument, "Usernames list cannot be empty")
	}

	domainChat := converter.FromProtoCreateChatRequest(req)

	createdChat, err := c.chatService.CreateChat(ctx, domainChat)
	if err != nil {
		log.Printf("CreateChat: Error creating chat: %v", err)
		return nil, status.Error(codes.Internal, "Failed to create chat")
	}

	response := &desc.CreateResponse{
		Id: createdChat.ID,
	}

	return response, nil
}

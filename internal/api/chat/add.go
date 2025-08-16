package chat

import (
	"context"
	"github.com/en7ka/chat-server/internal/converter"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) AddMemberToChat(ctx context.Context, req *desc.AddMemberToChatRequest) (*desc.AddMemberToChatResponse, error) {
	if req.GetChatId() <= 0 || req.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "chat_id and user_id must be positive")
	}

	domainMember := converter.FromProtoAddMemberRequest(req)
	memberID, err := i.chatService.AddMemberToChat(ctx, domainMember)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to add member")
	}

	return &desc.AddMemberToChatResponse{
		MemberId: memberID,
	}, nil
}

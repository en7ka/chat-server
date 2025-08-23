package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/en7ka/chat-server/internal/api/chat"
	"github.com/en7ka/chat-server/internal/converter"
	serviceMocks "github.com/en7ka/chat-server/internal/service/mocks"
	"github.com/en7ka/chat-server/internal/service/servinterface"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddMemberToChat(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) servinterface.ChatService

	type args struct {
		ctx context.Context
		req *desc.AddMemberToChatRequest
	}

	var (
		mc       = minimock.NewController(t)
		ctx      = context.Background()
		chatID   = gofakeit.Int64()
		userID   = gofakeit.Int64()
		memberID = gofakeit.Int64()

		req = &desc.AddMemberToChatRequest{
			ChatId: chatID,
			UserId: userID,
		}

		memberModel = converter.FromProtoAddMemberRequest(req)

		res = &desc.AddMemberToChatResponse{
			MemberId: memberID,
		}

		serviceErr = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.AddMemberToChatResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{ctx: ctx, req: req},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.AddMemberToChatMock.Expect(ctx, memberModel).Return(memberID, nil)
				return mock
			},
		},
		{
			name: "error case - service failure",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  status.Error(codes.Internal, "failed to add member"),
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.AddMemberToChatMock.Expect(ctx, memberModel).Return(0, serviceErr)
				return mock
			},
		},
		{
			name: "error case - invalid chat id",
			args: args{ctx: ctx, req: &desc.AddMemberToChatRequest{ChatId: 0, UserId: userID}},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "chat_id and user_id must be positive"),
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				return serviceMocks.NewChatServiceMock(mc)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewImplementation(chatServiceMock)

			response, err := api.AddMemberToChat(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}

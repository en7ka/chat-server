package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/en7ka/chat-server/internal/api/chat"
	"github.com/en7ka/chat-server/internal/converter"
	"github.com/en7ka/chat-server/internal/models"
	serviceMocks "github.com/en7ka/chat-server/internal/service/mocks"
	"github.com/en7ka/chat-server/internal/service/servinterface"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateChat(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) servinterface.ChatService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		chatID    = gofakeit.Int64()
		usernames = []string{gofakeit.Username(), gofakeit.Username()}

		req = &desc.CreateRequest{
			Usernames: usernames,
		}

		chatModel = converter.FromProtoCreateChatRequest(req)

		res = &desc.CreateResponse{
			Id: chatID,
		}

		serviceErr = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
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
				mock.CreateChatMock.Expect(ctx, chatModel).Return(&models.Chat{ID: chatID}, nil)
				return mock
			},
		},
		{
			name: "error case - service failure",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  status.Error(codes.Internal, "Failed to create chat"),
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, chatModel).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "error case - validation failure",
			args: args{ctx: ctx, req: &desc.CreateRequest{Usernames: []string{}}},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Usernames list cannot be empty"),
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				// No service calls expected
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

			response, err := api.CreateChat(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}

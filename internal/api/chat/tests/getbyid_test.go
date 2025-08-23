package tests

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestGetChat(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) servinterface.ChatService

	type args struct {
		ctx context.Context
		req *desc.GetChatRequest
	}

	var (
		mc     = minimock.NewController(t)
		ctx    = context.Background()
		chatID = gofakeit.Int64()

		req = &desc.GetChatRequest{
			Id: chatID,
		}

		chatModel = &models.Chat{
			ID:        chatID,
			Name:      gofakeit.Company(),
			Type:      models.ChatTypeGroup,
			CreatedAt: time.Now(),
		}

		res = &desc.GetChatResponse{
			Chat: converter.ToProtoChat(chatModel),
		}

		serviceErr = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetChatResponse
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
				mock.GetChatByIdMock.Expect(ctx, chatID).Return(chatModel, nil)
				return mock
			},
		},
		{
			name: "error case - service failure",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  serviceErr, // GetChat пробрасывает ошибку
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetChatByIdMock.Expect(ctx, chatID).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "error case - invalid chat id",
			args: args{ctx: ctx, req: &desc.GetChatRequest{Id: 0}},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Chat ID must be a positive number"),
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

			response, err := api.GetChat(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}

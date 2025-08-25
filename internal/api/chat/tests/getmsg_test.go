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

func TestGetChatMessages(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) servinterface.ChatService

	type args struct {
		ctx context.Context
		req *desc.GetMessagesRequest
	}

	var (
		mc     = minimock.NewController(t)
		ctx    = context.Background()
		chatID = gofakeit.Int64()

		req = &desc.GetMessagesRequest{
			ChatId: chatID,
		}

		messages = []*models.Message{
			{ID: gofakeit.Int64(), ChatID: chatID, FromUserID: gofakeit.Int64(), Text: gofakeit.Sentence(3), Timestamp: time.Now()},
			{ID: gofakeit.Int64(), ChatID: chatID, FromUserID: gofakeit.Int64(), Text: gofakeit.Sentence(4), Timestamp: time.Now()},
		}

		protoMessages = []*desc.Message{
			converter.ToProtoMessage(messages[0]),
			converter.ToProtoMessage(messages[1]),
		}

		res = &desc.GetMessagesResponse{
			Messages: protoMessages,
		}

		serviceErr = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetMessagesResponse
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
				mock.GetChatMessagesMock.Expect(ctx, chatID).Return(messages, nil)
				return mock
			},
		},
		{
			name: "success case - no messages",
			args: args{ctx: ctx, req: req},
			want: &desc.GetMessagesResponse{Messages: []*desc.Message{}},
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetChatMessagesMock.Expect(ctx, chatID).Return([]*models.Message{}, nil)
				return mock
			},
		},
		{
			name: "error case - service failure",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  serviceErr, // GetChatMessages пробрасывает ошибку
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetChatMessagesMock.Expect(ctx, chatID).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "error case - invalid chat id",
			args: args{ctx: ctx, req: &desc.GetMessagesRequest{ChatId: -1}},
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

			response, err := api.GetChatMessages(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)

			if tt.want != nil && response != nil {
				require.Equal(t, len(tt.want.Messages), len(response.Messages))
				for i := range tt.want.Messages {
					require.Equal(t, tt.want.Messages[i].Id, response.Messages[i].Id)
					require.Equal(t, tt.want.Messages[i].Text, response.Messages[i].Text)
				}
			} else {
				require.Equal(t, tt.want, response)
			}
		})
	}
}

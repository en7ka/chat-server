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

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) servinterface.ChatService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		chatID  = gofakeit.Int64()
		userID  = gofakeit.Int64()
		msgText = gofakeit.Sentence(5)

		req = &desc.SendMessageRequest{
			ChatId:     chatID,
			FromUserId: userID,
			Text:       msgText,
		}

		messageModel = converter.FromProtoSendMessageRequest(req)

		createdMsg = &models.Message{
			ID:         gofakeit.Int64(),
			ChatID:     chatID,
			FromUserID: userID,
			Text:       msgText,
			Timestamp:  time.Now(),
		}

		res = &desc.SendMessageResponse{
			Message: converter.ToProtoMessage(createdMsg),
		}

		serviceErr = errors.New("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.SendMessageResponse
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
				mock.SendMessageMock.Expect(ctx, messageModel).Return(createdMsg, nil)
				return mock
			},
		},
		{
			name: "error case - service failure",
			args: args{ctx: ctx, req: req},
			want: nil,
			err:  serviceErr, // SendMessage просто пробрасывает ошибку
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, messageModel).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "error case - invalid chat id",
			args: args{ctx: ctx, req: &desc.SendMessageRequest{ChatId: 0, FromUserId: userID, Text: msgText}},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "chat id must be greater than zero"),
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				return serviceMocks.NewChatServiceMock(mc)
			},
		},
		{
			name: "error case - invalid user id",
			args: args{ctx: ctx, req: &desc.SendMessageRequest{ChatId: chatID, FromUserId: 0, Text: msgText}},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "from user id must be greater than zero"),
			chatServiceMock: func(mc *minimock.Controller) servinterface.ChatService {
				return serviceMocks.NewChatServiceMock(mc)
			},
		},
		{
			name: "error case - empty text",
			args: args{ctx: ctx, req: &desc.SendMessageRequest{ChatId: chatID, FromUserId: userID, Text: "   "}},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "message text cannot be empty"),
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

			response, err := api.SendMessage(tt.args.ctx, tt.args.req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, response)
		})
	}
}

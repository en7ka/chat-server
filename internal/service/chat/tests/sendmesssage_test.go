package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/en7ka/chat-server/internal/client/db"
	dbMocks "github.com/en7ka/chat-server/internal/client/mocks"
	"github.com/en7ka/chat-server/internal/models"
	repoMocks "github.com/en7ka/chat-server/internal/repository/mocks"
	"github.com/en7ka/chat-server/internal/repository/repointerface"
	"github.com/en7ka/chat-server/internal/service/chat"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repointerface.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx     context.Context
		message *models.Message
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		msgID   = gofakeit.Int64()
		repoErr = errors.New("repository send message error")
		txErr   = errors.New("transaction manager error")

		serviceMessage = &models.Message{
			ChatID:     gofakeit.Int64(),
			FromUserID: gofakeit.Int64(),
			Text:       gofakeit.Sentence(10),
		}

		returnedMessage = &models.Message{
			ID:         msgID,
			ChatID:     serviceMessage.ChatID,
			FromUserID: serviceMessage.FromUserID,
			Text:       serviceMessage.Text,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *models.Message
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "successful send message",
			args: args{ctx: ctx, message: serviceMessage},
			want: returnedMessage,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, serviceMessage).Return(msgID, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "error_during_repository_send_message",
			args: args{ctx: ctx, message: serviceMessage},
			want: nil,
			// !!! ВОТ ЗДЕСЬ ИСПРАВЛЕНИЕ !!!
			err: fmt.Errorf("failed to send message in repository: %w", repoErr),
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, serviceMessage).Return(0, repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "error from transaction manager",
			args: args{ctx: ctx, message: serviceMessage},
			want: nil,
			err:  txErr,
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				return repoMocks.NewChatRepositoryMock(mc)
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommitedMock.Return(txErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)
			defer mc.Finish()

			chatRepositoryMock := tt.chatRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := chat.NewService(chatRepositoryMock, txManagerMock)

			message, err := service.SendMessage(tt.args.ctx, tt.args.message)

			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, message)
		})
	}
}

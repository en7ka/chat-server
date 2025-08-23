package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

func TestGetChatById(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repointerface.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx    context.Context
		chatID int64
	}

	var (
		mc      = minimock.NewController(t)
		ctx     = context.Background()
		chatID  = gofakeit.Int64()
		repoErr = errors.New("repository get chat error")
		txErr   = errors.New("transaction manager error")

		repoChat = &models.Chat{
			ID:        chatID,
			Name:      gofakeit.Company(),
			Type:      models.ChatTypeGroup,
			CreatedAt: time.Now(),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *models.Chat
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "successful get chat by id",
			args: args{ctx: ctx, chatID: chatID},
			want: repoChat,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.GetChatByIdMock.Expect(ctx, chatID).Return(repoChat, nil)
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
			name: "error from repository get chat by id",
			args: args{ctx: ctx, chatID: chatID},
			want: nil,
			err:  fmt.Errorf("failed to get chat by id from repository: %w", repoErr),
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.GetChatByIdMock.Expect(ctx, chatID).Return(nil, repoErr)
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
			args: args{ctx: ctx, chatID: chatID},
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

			chat, err := service.GetChatById(tt.args.ctx, tt.args.chatID)

			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, chat)
		})
	}
}

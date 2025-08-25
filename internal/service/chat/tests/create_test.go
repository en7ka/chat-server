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

func TestCreateChat(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repointerface.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		chat *models.Chat
	}

	var (
		mc       = minimock.NewController(t)
		ctx      = context.Background()
		chatID   = gofakeit.Int64()
		chatName = gofakeit.Company()

		repoErr = errors.New("repository create error")
		txErr   = errors.New("transaction manager error")

		serviceChat = &models.Chat{
			Name: chatName,
			Type: models.ChatTypeGroup,
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
			name: "successful create",
			args: args{
				ctx:  ctx,
				chat: serviceChat,
			},
			want: &models.Chat{ID: chatID, Name: chatName, Type: models.ChatTypeGroup},
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, serviceChat).Return(chatID, nil)
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
			name: "error during repository creation",
			args: args{
				ctx:  ctx,
				chat: serviceChat,
			},
			want: nil,
			// ИЗМЕНЕНИЕ: Ожидаем обернутую ошибку
			err: fmt.Errorf("failed to create chat in repository: %w", repoErr),
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateChatMock.Expect(ctx, serviceChat).Return(0, repoErr)
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
			args: args{
				ctx:  ctx,
				chat: serviceChat,
			},
			want: nil,
			err:  txErr,
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				// Репозиторий не должен быть вызван, если транзакция не началась
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

			createdChat, err := service.CreateChat(tt.args.ctx, tt.args.chat)

			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, createdChat)
		})
	}
}

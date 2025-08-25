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

func TestAddMemberToChat(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repointerface.ChatRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx    context.Context
		member *models.ChatMember
	}

	var (
		mc       = minimock.NewController(t)
		ctx      = context.Background()
		memberID = gofakeit.Int64()
		repoErr  = errors.New("repository add member error")
		txErr    = errors.New("transaction manager error")

		serviceMember = &models.ChatMember{
			ChatID: gofakeit.Int64(),
			UserID: gofakeit.Int64(),
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "successful add member",
			args: args{ctx: ctx, member: serviceMember},
			want: memberID,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.AddMemberToChatMock.Expect(ctx, serviceMember).Return(memberID, nil)
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
			name: "error during repository add member",
			args: args{ctx: ctx, member: serviceMember},
			want: 0,
			err:  fmt.Errorf("failed to add member in repository: %w", repoErr),
			chatRepositoryMock: func(mc *minimock.Controller) repointerface.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.AddMemberToChatMock.Expect(ctx, serviceMember).Return(0, repoErr)
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
			args: args{ctx: ctx, member: serviceMember},
			want: 0,
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

			newID, err := service.AddMemberToChat(tt.args.ctx, tt.args.member)

			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, newID)
		})
	}
}

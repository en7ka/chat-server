package chat

import (
	"context"
	"github.com/en7ka/chat-server/internal/repository/chat/model"
	repinf "github.com/en7ka/chat-server/internal/repository/repointerface"
	"github.com/jackc/pgx/v5/pgxpool"

	sq "github.com/Masterminds/squirrel"
)

// Константы для таблицы chats
const (
	tableNameChats       = "chat.chats"
	chatsIdColumn        = "id"
	chatsNameColumn      = "name"
	chatsTypeColumn      = "type"
	chatsCreatedAtColumn = "created_at"
	chatsIsDeletedColumn = "is_deleted"
)

// Константы для таблицы chat_members
const (
	tableNameChatMembers      = "chat.chat_members"
	chatMembersIdColumn       = "id"
	chatMembersChatIdColumn   = "chat_id"
	chatMembersUserIdColumn   = "user_id"
	chatMembersJoinedAtColumn = "joined_at"
	chatMembersIsAdminColumn  = "is_admin"
)

// Константы для таблицы messages
const (
	tableNameMessages        = "chat.messages"
	messagesIdColumn         = "id"
	messagesChatIdColumn     = "chat_id"
	messagesFromUserIdColumn = "from_user_id"
	messagesTextColumn       = "text"
	messagesTimestampColumn  = "timestamp"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repinf.ChatRepository {

	return &repo{db: db}
}

func (r *repo) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	qb := sq.Insert(tableNameChats).
		Columns(chatsNameColumn, chatsTypeColumn).
		Values(chat.Name, chat.Type).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + chatsIdColumn)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var chatId int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(&chatId); err != nil {
		return 0, err
	}

	return chatId, nil
}

func (r *repo) AddMemberToChat(ctx context.Context, member *model.ChatMember) (int64, error) {
	qb := sq.Insert(tableNameChatMembers).
		Columns(chatMembersIdColumn, chatMembersUserIdColumn, chatMembersIsAdminColumn).
		Values(member.ChatID, member.UserId, member.IsAdmin).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + chatsIdColumn)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var chatId int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(); err != nil {
		return 0, err
	}

	return chatId, nil
}

func (r *repo) SendMessage(ctx context.Context, message *model.Message) (int64, error) {
	qb := sq.Insert(tableNameMessages).
		Columns(messagesIdColumn, messagesFromUserIdColumn, messagesTextColumn).
		Values(message.ChatID, message.FromUserID, message.Text).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + messagesIdColumn)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var chatId int64
	if err := r.db.QueryRow(ctx, query, args...).Scan(&chatId); err != nil {
		return 0, err
	}

	return chatId, nil
}

func (r *repo) GetChatMessages(ctx context.Context, chatId int64) ([]*model.Message, error) {
	qb := sq.Select(
		messagesIdColumn,
		messagesChatIdColumn,
		messagesFromUserIdColumn,
		messagesTextColumn,
		messagesTimestampColumn,
	).
		From(tableNameMessages).
		Where(sq.Eq{messagesChatIdColumn: chatId}).
		OrderBy(messagesTimestampColumn + " ASC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.ChatID,
			&msg.FromUserID,
			&msg.Text,
			&msg.Timestamp,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

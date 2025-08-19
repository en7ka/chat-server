package chat

import (
	"context"
	"errors"

	"github.com/en7ka/chat-server/internal/client/db"
	"github.com/en7ka/chat-server/internal/converter"
	"github.com/en7ka/chat-server/internal/models"
	repoModel "github.com/en7ka/chat-server/internal/repository/chat/model"
	repinf "github.com/en7ka/chat-server/internal/repository/repointerface"
	"github.com/jackc/pgx/v5"

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
	db db.Client
}

func NewRepository(db db.Client) repinf.ChatRepository {

	return &repo{db: db}
}

func (r *repo) CreateChat(ctx context.Context, chat *models.Chat) (int64, error) {
	repoChat := converter.ToRepoChatFromDomain(chat)
	qb := sq.Insert(tableNameChats).
		Columns(chatsNameColumn, chatsTypeColumn).
		Values(chat.Name, repoChat.Type).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + chatsIdColumn)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.CreateChat",
		QueryRaw: query,
	}

	var chatId int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatId); err != nil {
		return 0, err
	}

	return chatId, nil
}

func (r *repo) AddMemberToChat(ctx context.Context, member *models.ChatMember) (int64, error) {
	repoMember := converter.ToRepoChatMemberFromDomain(member)

	qb := sq.Insert(tableNameChatMembers).
		Columns(chatMembersChatIdColumn, chatMembersUserIdColumn, chatMembersIsAdminColumn).
		Values(repoMember.ChatID, repoMember.UserID, repoMember.IsAdmin).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + chatMembersIdColumn)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.AddMemberToChat",
		QueryRaw: query,
	}

	var memberId int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(); err != nil {
		return 0, err
	}

	return memberId, nil
}

func (r *repo) SendMessage(ctx context.Context, message *models.Message) (int64, error) {
	repoMsg := converter.ToRepoMessageFromDomain(message)

	qb := sq.Insert(tableNameMessages).
		Columns(messagesChatIdColumn, messagesFromUserIdColumn, messagesTextColumn).
		Values(message.ChatID, message.FromUserID, repoMsg.Text).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING " + messagesIdColumn)

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.SendMessage",
		QueryRaw: query,
	}

	var messageId int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&messageId); err != nil {
		return 0, err
	}

	return messageId, nil
}

func (r *repo) GetChatMessages(ctx context.Context, chatId int64) ([]*models.Message, error) {
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

	q := db.Query{
		Name:     "user_repository.GetChatMessages",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repoMessages []*repoModel.Message
	for rows.Next() {
		var msg repoModel.Message
		if err = rows.Scan(
			&msg.ChatID,
			&msg.FromUserID,
			&msg.Text,
			&msg.Timestamp); err != nil {
			return nil, err
		}
		repoMessages = append(repoMessages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	domainMessages := make([]*models.Message, 0, len(repoMessages))
	for _, repoMsg := range repoMessages {
		domainMessages = append(domainMessages, converter.ToDomainMessageFromRepo(repoMsg))
	}

	return domainMessages, nil
}

func (r *repo) GetChatById(ctx context.Context, chatId int64) (*models.Chat, error) {
	qb := sq.Select(
		chatsIdColumn,
		chatsNameColumn,
		chatsTypeColumn,
		chatsCreatedAtColumn,
		chatsIsDeletedColumn,
	).
		From(tableNameChats).
		Where(sq.Eq{chatsIdColumn: chatId}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetChatMessages",
		QueryRaw: query,
	}

	var repoChat repoModel.Chat
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&repoChat.ID,
		&repoChat.Name,
		&repoChat.Type,
		&repoChat.CreatedAt,
		&repoChat.IsDeleted); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return converter.ToDomainChatFromRepo(&repoChat), nil
}

package cmd

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"log"
)

const (
	dbDSN = "host=localhost port=5433 dbname=chat_db user=data-user password=note-password sslmode=disable"
)

type PostgresConfig struct {
	con pgx.Conn
}

// подключение к бд
func InitPostgresConfig() *PostgresConfig {
	con, err := pgx.Connect(context.Background(), dbDSN)
	if err != nil {
		log.Fatal("ошибка подключения к бд", err)
	}
	return &PostgresConfig{con: *con}
}

// закрытие бд
func (s *PostgresConfig) CloseCon() {
	err := s.con.Close(context.Background())
	if err != nil {
		log.Fatal("ошибка в закрытии бд", err)
	}
}

// интерфейс для работы с бд
type PostgresInterface interface {
	CreateChat(users IDs) (int64, error)
	DeleteChat(id int64) error
	SendMessageChat(message Message) error
}

// реализация CreateChat
func (s *PostgresConfig) CreateChat(users IDs) (*int64, error) {
	var chatID int64
	ctx := context.Background()

	query := "INSERT INTO chat.chats (type) VALUES ($1) RETURNING id"

	if err := s.con.QueryRow(ctx, query, "group").Scan(&chatID); err != nil {
		log.Printf("Ошибка при создании чата: %v", err)
		return nil, fmt.Errorf("ошибка при создании чата: %w", err)
	}

	queryBuilder := sq.Insert("chat.chat_members").Columns("chat_id", "user_id")

	for _, userID := range users {
		queryBuilder = queryBuilder.Values(chatID, userID)
	}

	sqlStr, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Printf("Ошибка при формировании запроса: %v", err)
		return nil, fmt.Errorf("ошибка при формировании запроса на добавление участников: %w", err)
	}
	_, err = s.con.Exec(ctx, sqlStr, args...)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса на добавление участников: %v", err)
		return nil, fmt.Errorf("ошибка при выполнении запроса на добавление участников: %w", err)
	}
	log.Printf("Создан чат с id: %d. В него добавлены пользователи: %+v", chatID, users)

	return &chatID, nil
}

// реализация DeleteChat
func (s *PostgresConfig) DeleteChat(chatID int64) error {
	ctx := context.Background()
	del := "UPDATE chat.chats SET is_deleted = TRUE WHERE id = $1"

	_, err := s.con.Exec(ctx, del, chatID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении чата с id %d: %w", chatID, err)
	}
	log.Printf("Удален чат с id: %d", chatID)

	return nil
}

// реализация SendMessage
func (s *PostgresConfig) SendMessageChat(message Message) error {
	ctx := context.Background()
	var chatIsDeleted bool

	delBl := "SELECT is_deleted FROM chat.chats WHERE id = $1"
	if err := s.con.QueryRow(ctx, delBl, message.ChatID).Scan(&chatIsDeleted); err != nil {
		return fmt.Errorf("ошибка при проверке чата с id %d: %w", message.ChatID, err)
	}

	if chatIsDeleted {
		return fmt.Errorf("чат с id %d удален", message.ChatID)
	}

	insertQuery := "INSERT INTO chat.messages (chat_id, from_user_id, text) VALUES ($1, $2, $3)"
	_, err := s.con.Exec(ctx, insertQuery, message.ChatID, message.FromUID, message.Body)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения в чат %d: %v", message.ChatID, err)
		return fmt.Errorf("ошибка при отправке сообщения: %w", err)
	}
	log.Printf("Пользователь %s отправил сообщение %s в чат %d", message.From, message.Body, message.ChatID)

	return nil
}

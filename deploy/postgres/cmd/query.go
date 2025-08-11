package cmd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

const (
	dbDSN = "DEFAULT_DSN"
)

type PostgresConfig struct {
	con pgx.Conn
}

// подключение к бд
func InitPostgresConfig() *PostgresConfig {
	con, err := pgx.Connect(context.Background(), dbDSN)
	if err != nil {
		log.Fatalf("ошибка подключения к бд: %v", err)
	}

	return &PostgresConfig{con: *con}
}

// закрытие бд
func (s *PostgresConfig) CloseCon() {
	err := s.con.Close(context.Background())
	if err != nil {
		log.Printf("ошибка в закрытии бд: %v", err)
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

	tx, err := s.con.Begin(ctx)
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v", err)
		return nil, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}

	defer tx.Rollback(ctx)

	query := "INSERT INTO chat.chats (type) VALUES ($1) RETURNING id"
	if err := tx.QueryRow(ctx, query, "group").Scan(&chatID); err != nil {
		log.Printf("Ошибка при создании чата: %v", err)
		return nil, fmt.Errorf("ошибка при создании чата: %w", err)
	}

	insertQuery := "INSERT INTO chat.chat_members (chat_id, user_id) VALUES ($1, $2)"
	_, err = tx.Exec(ctx, insertQuery, chatID, users)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса на добавление участников: %v", err)
		return nil, fmt.Errorf("ошибка при выполнении запроса на добавление участников: %w", err)
	}
	//если все успешно, коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return nil, fmt.Errorf("ошибка при коммите транзакции: %w", err)
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
	insertQuery := `
			INSERT INTO chat.chat_members
			SELECT $1, $2, $3
			FROM chat.chats
			WHERE id = $1 AND is_deleted = FALSE`

	com, err := s.con.Exec(ctx, insertQuery, message.ChatID, message.FromUID, message.Body)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения в чат %d: %v", message.ChatID, err)
		return fmt.Errorf("ошибка при отправке сообщения: %w", err)
	}
	//проверяем была ли вставлена хоть одна строка
	if com.RowsAffected() == 0 {
		return fmt.Errorf("не удалось отправить сообщение в чат %d: чат не найден или удален", message.ChatID)
	}
	log.Printf("Пользователь %s отправил сообщение %s в чат %d", message.From, message.Body, message.ChatID)

	return nil
}

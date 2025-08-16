package main

import (
	"context"
	chatApi "github.com/en7ka/chat-server/internal/api/chat"
	"github.com/en7ka/chat-server/internal/config"
	chatRepo "github.com/en7ka/chat-server/internal/repository/chat"
	chatService "github.com/en7ka/chat-server/internal/service/chat"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	ctx := context.Background()

	//считываем переменные окружения
	err := config.Load(".env")
	if err != nil {
		log.Fatalf("ошибка к подключению к .env: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("ошибка к подключению с grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("ошибка к подключению к pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Printf("ошибка в прослушивании: %v", err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Printf("ошибка в подключении: %v", err)
	}
	defer pool.Close()

	noteRepo := chatRepo.NewRepository(pool)
	noteService := chatService.NewService(noteRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatAPIServer(s, chatApi.NewImplementation(noteService))

	log.Printf("Сервер запущен на %s", grpcConfig.Address())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("ошибка в запуске сервера: %v", err)
	}
}

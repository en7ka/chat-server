package main

import (
	"context"
	"fmt"
	userv1 "github.com/en7ka/auth/pkg/user_v1"
	"github.com/en7ka/chat-server/deploy/postgres/cmd"
	chatv1 "github.com/en7ka/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

func main() {
	authConn, err := grpc.Dial("127.0.0.1:50501", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("не удалось подключиться к auth-серверу: %v", err)
	}
	defer authConn.Close()

	authClient := userv1.NewUserAPIClient(authConn)

	config := cmd.InitPostgresConfig()
	defer config.CloseCon()

	lis, err := net.Listen("tcp", "127.0.0.1:50050")
	if err != nil {
		log.Fatalf("не удалось прослушать порт: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chatv1.RegisterChatAPIServer(s, &server{
		authClient: authClient,
		config:     config,
	})

	log.Printf("Сервер запущен на порту 50050")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	chatv1.UnimplementedChatAPIServer
	authClient userv1.UserAPIClient
	config     *cmd.PostgresConfig
}

func (s *server) Create(ctx context.Context, req *chatv1.CreateRequest) (*chatv1.CreateResponse, error) {
	log.Printf("Users: %v", req.GetUsernames())

	resp, err := s.authClient.GetByUsernames(ctx, &userv1.GetByUsernamesRequest{
		Usernames: req.GetUsernames(),
	})
	if err != nil {
		log.Printf("Ошибка при пакетном получении пользователей: %v", err)
		return nil, fmt.Errorf("ошибка при получении пользователей: %w", err)
	}
	//проверим что все пользователи были найдены
	if len(resp.GetUsers()) != len(req.GetUsernames()) {
		log.Printf("Не все пользователи найдены. Запрошено: %d, найдено: %d", len(req.GetUsernames()), len(resp.GetUsers()))
		return nil, fmt.Errorf("не все пользователи найдены: запрошено %d, найдено %d", len(req.GetUsernames()), len(resp.GetUsers()))
	}

	//собераем id
	userIDs := make([]int64, 0, len(resp.GetUsers()))
	for i, user := range resp.GetUsers() {
		userIDs[i] = user.GetId()
	}
	log.Printf("получены ID пользователей: %v", userIDs)

	id, err := s.config.CreateChat(userIDs)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании чата: %w", err)
	}

	return &chatv1.CreateResponse{Id: *id}, nil
}

func (s *server) Delete(ctx context.Context, req *chatv1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Удаление чата с id: %d", req.GetId())

	if err := s.config.DeleteChat(req.GetId()); err != nil {
		return nil, fmt.Errorf("ошибка при удалении чата: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Отправка сообщения в чат %d от %s", req.GetChatId(), req.GetFrom())

	if req.GetChatId() == 0 {
		return nil, fmt.Errorf("chat ID is required")
	}

	userResp, err := s.authClient.GetByUsername(ctx, &userv1.GetByUsernameRequest{
		Username: req.GetFrom(),
	})
	if err != nil {
		log.Printf("Ошибка при получении ID отправителя %s: %v", req.GetFrom(), err)
		return nil, fmt.Errorf("ошибка при получении ID отправителя %s: %w", req.GetFrom(), err)
	}
	fromUserID := userResp.GetId()

	msg := cmd.Message{
		ChatID:  req.GetChatId(),
		From:    req.GetFrom(),
		FromUID: fromUserID,
		Body:    req.GetText(),
		Time:    req.GetTime().AsTime(),
	}

	if err := s.config.SendMessageChat(msg); err != nil {
		return nil, fmt.Errorf("ошибка при отправке сообщения: %w", err)
	}

	return &emptypb.Empty{}, nil
}

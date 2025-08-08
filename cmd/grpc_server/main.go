package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	desc "github.com/en7ka/chat-server/pkg/chat_v1"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Received Create chat request for users: %v", req.GetUsernames())

	return &desc.CreateResponse{
		Id: req.GetId(),
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Received Delete chat request for ID: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf(
		"Received SendMessage request: From=%s, Text=%s, Timestamp=%s",
		req.GetFromUsername(),
		req.GetText(),
		req.GetTimestamp().AsTime().Format(time.RFC3339),
	)
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	desc.RegisterChatV1Server(s, &server{})
	reflection.Register(s)

	log.Printf("Starting gRPC server on port %d", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

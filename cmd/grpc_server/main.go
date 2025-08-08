package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	desc "github.com/en7ka/chat-server/pkg/chat_v1"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatAPIServer
}

var nextID int64

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Received Create chat request for users: %v", req.GetUsernames())
	id := atomic.AddInt64(&nextID, 1)
	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Received Delete chat request for id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("SendMessage from=%s text=%s time=%s", req.GetFrom(), req.GetText(), req.GetTime().AsTime().Format(time.RFC3339))
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	desc.RegisterChatAPIServer(s, &server{})
	reflection.Register(s)

	log.Printf("Starting gRPC server on port %d", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

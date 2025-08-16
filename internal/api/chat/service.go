package chat

import (
	usserv "github.com/en7ka/chat-server/internal/service/servinterface"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
)

type Implementation struct {
	desc.UnimplementedChatAPIServer
	chatService usserv.ChatService
}

func NewImplementation(chatService usserv.ChatService) *Implementation {

	return &Implementation{chatService: chatService}
}

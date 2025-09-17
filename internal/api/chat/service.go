package chat

import (
	"sync"

	usserv "github.com/en7ka/chat-server/internal/service/servinterface"
	desc "github.com/en7ka/chat-server/pkg/chat_v1"
)

// subscriber представляет одного подключенного клиента.
type subscriber struct {
	id string
	ch chan *desc.Message
}
type Controller struct {
	desc.UnimplementedChatAPIServer
	chatService usserv.ChatService

	mu          sync.RWMutex
	subscribers map[int64][]*subscriber
}

func NewImplementation(chatService usserv.ChatService) *Controller {

	return &Controller{
		chatService: chatService,
		subscribers: make(map[int64][]*subscriber),
	}
}

package chat

import (
	"log"

	desc "github.com/en7ka/chat-server/pkg/chat_v1"
	"github.com/google/uuid"
)

func (c *Controller) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatAPI_ConnectChatServer) error {
	log.Printf("Клиент подключается к чату ID: %d", req.GetId())

	c.mu.Lock()

	newSub := &subscriber{
		id: uuid.NewString(),
		ch: make(chan *desc.Message, 100),
	}

	c.subscribers[req.GetId()] = append(c.subscribers[req.GetId()], newSub)

	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		subs := c.subscribers[req.GetId()]
		for i, sub := range subs {
			if sub.id == newSub.id {
				c.subscribers[req.GetId()] = append(subs[:i], subs[i+1:]...)
				close(newSub.ch)
				log.Printf("Клиент %s отписался от чата %d", newSub.id, req.GetId())
				break
			}
		}
	}()

	for {
		select {
		case msg, ok := <-newSub.ch:
			if !ok {
				return nil
			}
			if err := stream.Send(msg); err != nil {
				log.Printf("Ошибка отправки сообщения клиенту: %v", err)
				return err
			}

		case <-stream.Context().Done():
			log.Printf("Клиент отключился от чата ID: %d. Контекст завершен.", req.GetId())
			return stream.Context().Err()
		}
	}
}

package chat

import (
	"github.com/teablog/tea/internal/logger"
)

type Responser interface {
	Bytes() []byte
	GetArticleID() string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clientHub map[string]map[*Client]struct{}

	// Inbound messages from the clients.
	broadcast chan Responser

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Responser),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clientHub:  make(map[string]map[*Client]struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// 注册客户端
			if _, ok := h.clientHub[client.articleId]; !ok {
				// 客户端已存在 ？
			} else {
				h.clientHub[client.articleId] = map[*Client]struct{}{client: {}}
			}
		case client := <-h.unregister:
			if clients, ok := h.clientHub[client.articleId]; ok {
				delete(clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			logger.Debugf("广播一条新消息: %s", message.Bytes())
			if clients, ok := h.clientHub[message.GetArticleID()]; ok {
				for client, _ := range clients {
					select {
					case client.send <- message.Bytes():
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

package chat

import (
	"dyc/internal/consts"
	"dyc/internal/logger"
)

type Responser interface {
	Members() []string
	Bytes() []byte
	GetChannelID() string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

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
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if otherClient, ok := h.clients[client.account.Id]; ok {
				otherClient.conn.Close()
				client.send <- NewTipMessage("该账号的其他连接已关闭").Bytes()
			}
			h.clients[client.account.Id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.account.Id]; ok {
				delete(h.clients, client.account.Id)
				close(client.send)
			}
		case message := <-h.broadcast:
			logger.Debugf("广播一条新消息: %s", message.Bytes())
			if message.GetChannelID() == consts.GlobalChannelId {
				for id, client := range h.clients {
					select {
					case client.send <- message.Bytes():
					default:
						close(client.send)
						delete(h.clients, id)
					}
				}
			} else {
				for _, id := range message.Members() {
					if client, ok := h.clients[id]; ok {
						select {
						case client.send <- message.Bytes():
						default:
							close(client.send)
							delete(h.clients, id)
						}
					}
				}
			}
		}
	}
}

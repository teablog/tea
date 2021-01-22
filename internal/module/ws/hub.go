package ws

type Responser interface {
	Bytes() []byte
	GetArticleID() string
	GetAccountID() string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered uid
	uuid map[string]int

	// Inbound messages from the clients.
	broadcast chan Responser

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var hub *Hub

func init() {
	hub = &Hub{
		broadcast:  make(chan Responser),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		uuid:       make(map[string]int),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// 注册客户端
			client.send <- client.hub.Count().Bytes()
			h.uuid[client.uuid]++
		case client := <-h.unregister:
			h.uuid[client.uuid]--
			if h.uuid[client.uuid] == 0 {
				delete(h.uuid, client.uuid)
			}
		}
	}
}

func (h *Hub) Count() *ServerMessage {
	return &ServerMessage{Type: OnlineMsg, Content: len(h.uuid)}
}

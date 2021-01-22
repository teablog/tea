package ws

type Responser interface {
	Bytes() []byte
	GetArticleID() string
	GetAccountID() string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	client map[*Client]struct{}

	// Registered uid
	uuid map[string]int

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
		client:     make(map[*Client]struct{}),
		uuid:       make(map[string]int),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// 注册客户端
			h.client[client] = struct{}{}
			//if _, ok := h.uuid[client.]

		case client := <-h.unregister:
			delete(h.client, client)
		}
	}
}

func (h *Hub) Count() *ServerMessage {
	return &ServerMessage{Type: OnlineMsg, Content: len(h.client)}
}

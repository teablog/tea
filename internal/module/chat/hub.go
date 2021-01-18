package chat

import (
	"fmt"
	"github.com/teablog/tea/internal/logger"
	"sync"
	"time"
)

type Responser interface {
	Bytes() []byte
	GetArticleID() string
	GetAccountID() string
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

	// message storage
	messageHub *messageHub
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Responser),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clientHub:  make(map[string]map[*Client]struct{}),
		messageHub: newMessageHub(10),
	}
}

func (h *Hub) Run() {
	go h.messageHub.Start()
	for {
		select {
		case client := <-h.register:
			// 注册客户端
			if _, ok := h.clientHub[client.articleId]; !ok {
				// 客户端已存在 ？
				h.clientHub[client.articleId] = make(map[*Client]struct{})
				h.clientHub[client.articleId][client] = struct{}{}
			} else {
				h.clientHub[client.articleId][client] = struct{}{}
			}
		case client := <-h.unregister:
			if clients, ok := h.clientHub[client.articleId]; ok {
				delete(clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			logger.Debugf("广播一条新消息: %s", message.Bytes())
			if clients, ok := h.clientHub[message.GetArticleID()]; ok {
				h.messageHub.storager <- message.(*ServerMessage)
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

// MessageHub 临时缓存消息，超过配置就会写入磁盘
// 1. 10条
// 2. 1分钟
type messageHub struct {
	pos      int // 当前位置
	sw       sync.Mutex
	messages []*ServerMessage
	storager chan *ServerMessage
	len      int
}

func newMessageHub(len int) *messageHub {
	return &messageHub{
		pos:      0,
		sw:       sync.Mutex{},
		messages: make([]*ServerMessage, len),
		storager: make(chan *ServerMessage, len),
		len:      len,
	}
}

// store: 同一篇文章，同一个账户连续发表的内容合并
func (h *messageHub) store() {
	messages := make(ServerMessageSlice, 0, h.len)
	for i, j := 0, 0; i < h.pos; i++ {
		if i > 0 {
			if messages[j-1].GetArticleID() == h.messages[i].GetArticleID() && messages[j-1].GetAccountID() == h.messages[i].GetAccountID() {
				messages[j-1].Content = fmt.Sprintf("%s <br /> %s", messages[j-1].Content, h.messages[i].Content)
				continue
			}
		}
		messages = append(messages, h.messages[i])
		j++
	}
	messages.store()
	h.pos = 0
}

func (h *messageHub) Start() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				h.sw.Lock()
				h.store()
				h.sw.Unlock()
			}
		}
	}()
	for message := range h.storager {
		h.sw.Lock()
		if h.pos < h.len {
			h.messages[h.pos] = message
			h.pos++
		} else {
			h.store()
		}
		h.sw.Unlock()
	}
}

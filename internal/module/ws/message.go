package ws

import (
	"encoding/json"
	"github.com/teablog/tea/internal/logger"
)

// 消息来源

type msgType string

const (
	OnlineMsg msgType = "ONLINE"
)

type ServerMessage struct {
	Type msgType `json:"type"`
	// 内容
	Content interface{} `json:"content"`
}

func (row *ServerMessage) Bytes() []byte {
	data, err := json.Marshal(row)
	if err != nil {
		logger.Wrapf(err, "ws json err")
	}
	return data
}

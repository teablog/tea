package chat

import (
	"bytes"
	"dyc/internal/consts"
	"dyc/internal/db"
	"dyc/internal/helper"
	"dyc/internal/logger"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"time"
)

// 消息来源

type msgType string

const (
	TextMsg   msgType = "TEXT"
	ImgMsg    msgType = "IMAGE"
	FileMsg   msgType = "FILE"
	SystemMsg msgType = "SYSTEM"
	TipMsg    msgType = "TIP"
	OnlineMsg msgType = "ONLINE"
)

type shortAcct struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type ServerMessage struct {
	// 消息id
	Id string `json:"id"`
	// 时间
	Date time.Time `json:"date"`
	// 发送者
	Sender shortAcct `json:"sender"`
	// 类型
	Type msgType `json:"type"`
	// 内容
	Content string `json:"content"`
	// channel
	ChannelId string `json:"channel_id"`
}

type ClientMessage struct {
	ChannelId string `json:"channel_id"`
	Content   string `json:"content"`
}

func NewDefaultMsg(c *Client, msg string, channelId string) *ServerMessage {
	m := &ServerMessage{
		Content: msg,
		Sender: shortAcct{
			Id:   c.account.Id,
			Name: c.account.Name,
		},
		Type:      TextMsg,
		Date:      time.Now(),
		ChannelId: channelId,
	}
	m.Id = m.GenId()
	m.store()
	return m
}

func NewTipMessage(msg string) *ServerMessage {
	m := &ServerMessage{
		Content: msg,
		Date:    time.Time{},
		Sender:  shortAcct{Id: "0", Name: "系统消息"},
		Type:    TipMsg,
	}
	m.Id = m.GenId()
	return m
}

func (m *ServerMessage) GenId() string {
	var buf bytes.Buffer
	buf.WriteString(m.Date.String())
	buf.WriteString(m.Sender.Id)
	buf.WriteString(m.Content)
	return helper.Md532(buf.Bytes())
}

func (m *ServerMessage) Bytes() []byte {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		panic(errors.Wrap(err, "message json encode error"))
	}
	return buf.Bytes()
}

func (m *ServerMessage) Members() []string {
	ids, err := ChannelMembers.MembersIds(m.ChannelId)
	if err != nil {
		return nil
	}
	return ids
}

func (m *ServerMessage) GetChannelID() string {
	return m.ChannelId
}

func (m *ServerMessage) store() bool {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		logger.Wrapf(err, "message store json error")
		return false
	}
	if m.Id == "" {
		return false
	}
	res, err := db.ES.Index(
		consts.IndicesMessageConst,
		strings.NewReader(buf.String()),
		db.ES.Index.WithDocumentID(m.Id),
	)
	if err != nil {
		logger.Wrapf(err, "message store es error")
		return false
	}
	defer res.Body.Close()
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		logger.Errorf("[%d] es response: %s", res.StatusCode, string(resp))
		return false
	}
	return true
}

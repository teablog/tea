package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/account"
	"io/ioutil"
	"sort"
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

type ServerMessage struct {
	// 消息id
	Id string `json:"id"`
	// 时间
	Date time.Time `json:"date"`
	// 发送者
	Sender *account.Account `json:"sender"`
	// 类型
	Type msgType `json:"type"`
	// 内容
	Content string `json:"content"`
	// channel
	ArticleId string `json:"article_id"`
}

type ClientMessage struct {
	Content   string  `json:"content"`
	ArticleId string  `json:"article_id"`
	Type      msgType `json:"type"`
}

// 倒排获取30条
// 然后按照时间排序
type serverMessageSlice []*ServerMessage

func (m serverMessageSlice) Len() int {
	return len(m)
}

func (m serverMessageSlice) Less(i, j int) bool {
	return (m)[i].Date.Before((m)[j].Date)
}

func (m serverMessageSlice) Swap(i, j int) {
	(m)[i], (m)[j] = (m)[j], (m)[i]
}

func NewMessage(c *Client, cm ClientMessage) *ServerMessage {
	m := &ServerMessage{
		Content:   cm.Content,
		Sender:    c.account,
		Type:      cm.Type,
		Date:      time.Now(),
		ArticleId: cm.ArticleId,
	}
	m.Id = m.GenId()
	m.store()
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

func (m *ServerMessage) GetArticleID() string {
	return m.ArticleId
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
		consts.IndicesMessagesConst,
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

var Message *_message

type _message struct{}

func (*_message) FindMessages(articleId string, before time.Time) (int, serverMessageSlice, error) {
	query := fmt.Sprintf(`
{
  "query": {
    "bool": {
      "must": [
        {"term": {"article_id":  "%s"}},
        {"range": { "date": {"lt": "%s"}}}
      ]
    }
  },
  "sort": { "date": { "order": "desc" } },
  "size": 20
}`, articleId, before.Format(consts.EsTimeFormat))
	logger.Debugf("[ES query]: %s", query)
	resp, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMessagesConst),
		db.ES.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		logger.Errorf("[ES] %s search error: %s", consts.IndicesMessagesConst, err.Error())
		return 0, nil, errors.New("消息获取失败～")
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("[ES] response body read failed: %s", err.Error())
		return 0, nil, errors.New("消息获取失败～")
	}
	if resp.IsError() {
		logger.Errorf("[ES] response error: %s", string(res))
		return 0, nil, errors.New("消息获取失败～")
	}
	var r db.ESListResponse
	if err := json.Unmarshal(res, &r); err != nil {
		logger.Errorf("[json] unmarshal err: %s\n%s", err.Error(), string(res))
		return 0, nil, errors.New("消息获取失败～")
	}
	type hits []struct {
		Source *ServerMessage `json:"_source"`
		Id     string         `json:"_id"`
	}

	data := make(hits, 0)
	if err := json.Unmarshal(r.Hits.Hits, &data); err != nil {
		logger.Errorf("[json] unmarshal err: %s\n%s", err.Error(), string(r.Hits.Hits))
		return 0, nil, errors.New("消息获取失败～")
	}

	rows := make(serverMessageSlice, 0)
	for _, v := range data {
		rows = append(rows, v.Source)
	}

	sort.Sort(rows)

	return r.Hits.Total.Value, rows, nil
}

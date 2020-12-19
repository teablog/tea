package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/middleware"
	"github.com/teablog/tea/internal/module/account"
	"github.com/teablog/tea/internal/validate"
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
	Content   string  `json:"content" form:"content"`
	ArticleId string  `json:"article_id" form:"article_id"`
	Type      msgType `json:"type" form:"type"`
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

func NewMessage(acct *account.Account, cm ClientMessage) *ServerMessage {
	m := &ServerMessage{
		Content:   cm.Content,
		Sender:    acct,
		Type:      cm.Type,
		Date:      time.Now(),
		ArticleId: cm.ArticleId,
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

func (m *ServerMessage) GetArticleID() string {
	return m.ArticleId
}

func (m *ServerMessage) GetAccountID() string {
	return m.Sender.Id
}

type ServerMessageSlice []*ServerMessage

func (rows ServerMessageSlice) store() bool {
	if len(rows) == 0 {
		return false
	}
	buf := bytes.NewBuffer(nil)
	for _, v := range rows {
		buf.WriteString(fmt.Sprintf(`{ "index" : { "_index" : "%s", "_id" : "%s" } }`, consts.IndicesMessagesConst, v.Id))
		buf.WriteString("\n")
		buf.Write(v.Bytes())
	}
	logger.Debugf("[bulk]: %s", buf.String())
	res, err := db.ES.Bulk(buf)
	if err != nil {
		logger.Wrapf(err, "[ES] server messages bulk save err: %s", err.Error())
		return false
	}
	defer res.Body.Close()
	if res.IsError() {
		resp, _ := ioutil.ReadAll(res.Body)
		logger.Errorf("[ES] es response: %s", res.StatusCode, string(resp))
		return false
	}
	return true
}

var Message *_message

type _message struct {
}

// FindMessages 获取评论列表
func (*_message) FindMessages(req validate.MessagesValidator) (int, serverMessageSlice, error) {
	var (
		before time.Time
		after  time.Time
	)
	must := []string{fmt.Sprintf(`{"term": {"article_id": "%s"}}`, req.ArticleId)}
	if req.Before > 0 {
		before = time.Unix(req.Before/1000, int64(req.Before%1000)*1000000)
	}
	if req.After > 0 {
		after = time.Unix(req.After/1000, int64(req.After%1000)*1000000)
	}
	if !before.IsZero() {
		must = append(must, fmt.Sprintf(fmt.Sprintf(`{"range": { "date": {"lt": "%s"}}}`, before.Format(consts.EsTimeFormat))))
	}
	if !after.IsZero() {
		must = append(must, fmt.Sprintf(fmt.Sprintf(`{"range": { "date": {"gt": "%s"}}}`, after.Format(consts.EsTimeFormat))))
	}
	var (
		order       = "desc"
		size  int64 = 20
		from        = ``
	)
	if req.Sort == "asc" {
		order = "asc"
	}
	if req.Size > 0 {
		size = req.Size
	}
	if req.Page > 0 {
		from = fmt.Sprintf(`,"from": %d`, (req.Page-1)*size)
	}
	query := fmt.Sprintf(`{"query": {"bool": {"must": [%s]}}, "sort": { "date": { "order": "%s" } }, "size": %d %s}`, strings.Join(must, ","), order, size, from)

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

// SendMessage 评论
func (*_message) SendMessage(ctx *gin.Context, hub *Hub, message ClientMessage) error {
	acct := middleware.GetAccount(ctx)
	if acct != nil {
		hub.broadcast <- NewMessage(acct, message)
		return nil
	} else {
		return errors.New("")
	}
}

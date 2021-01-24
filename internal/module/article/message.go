package article

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
	"math"
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

func NewMessage(acct *account.Account, cm *ClientMessage) *ServerMessage {
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
	m.Id = helper.Md532(buf.Bytes())
	return m.Id
}

func (m *ServerMessage) Bytes() []byte {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		panic(errors.Wrap(err, "message json encode error"))
	}
	return buf.Bytes()
}

var Msg *_message

type _message struct{}

// FindMessages 获取评论列表
func (m *_message) FindMessages(ctx *gin.Context) {
	req := new(validate.MessagesValidator)
	if err := ctx.ShouldBind(req); err != nil {
		helper.ServerErr(ctx)
		return
	}
	must := []string{fmt.Sprintf(`{"term": {"article_id": "%s"}}`, req.ArticleId)}
	c, err := m.count(req.ArticleId)
	if err != nil {
		logger.Wrapf(err, "message count err ")
		helper.ServerErr(ctx)
		return
	}
	var (
		order = "desc"
		size  = 100
	)
	if req.Sort == "asc" {
		order = "asc"
	}
	page := int(math.Ceil(float64(c) / float64(size)))
	if req.Page > 0 {
		page = req.Page
	}
	query := fmt.Sprintf(`{"query": {"bool": {"must": [%s]}}, "sort": { "date": { "order": "%s" } }, "size": %d ,"from": %d}`, strings.Join(must, ","), order, size, (page-1)*size)
	logger.Debugf("[ES query]: %s", query)
	resp, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesMessagesConst),
		db.ES.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		logger.Errorf("[ES] %s search error: %s", consts.IndicesMessagesConst, err.Error())
		helper.ServerErr(ctx)
		return
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("[ES] response body read failed: %s", err.Error())
		helper.ServerErr(ctx)
		return
	}
	if resp.IsError() {
		logger.Errorf("[ES] response error: %s", string(res))
		helper.ServerErr(ctx)
		return
	}
	var r db.ESListResponse
	if err := json.Unmarshal(res, &r); err != nil {
		logger.Errorf("[json] unmarshal err: %s\n%s", err.Error(), string(res))
		helper.ServerErr(ctx)
		return
	}
	type hits []struct {
		Source *ServerMessage `json:"_source"`
		Id     string         `json:"_id"`
	}
	data := make(hits, 0)
	if err := json.Unmarshal(r.Hits.Hits, &data); err != nil {
		logger.Errorf("[json] unmarshal err: %s\n%s", err.Error(), string(r.Hits.Hits))
		helper.ServerErr(ctx)
		return
	}
	rows := make(serverMessageSlice, 0)
	for _, v := range data {
		rows = append(rows, v.Source)
	}
	sort.Sort(rows)
	helper.Success(ctx, gin.H{"total": r.Hits.Total.Value, "list": rows, "page": page, "size": size})
	return
}

func (*_message) count(articleId string) (int, error) {
	res, err := db.ES.Count(
		db.ES.Count.WithIndex(consts.IndicesMessagesConst),
		db.ES.Count.WithBody(strings.NewReader(fmt.Sprintf(`{ "query": { "term": { "article_id": "%s" } } }`, articleId))),
	)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.IsError() {
		data, _ := ioutil.ReadAll(res.Body)
		return 0, errors.Errorf("es response err: %s", string(data))
	}
	type r struct {
		Count int `json:"count"`
	}
	rs := new(r)
	if err := json.NewDecoder(res.Body).Decode(rs); err != nil {
		return 0, err
	}
	return rs.Count, nil
}

func (*_message) Comment(ctx *gin.Context) {
	a := middleware.GetAccount(ctx)
	cm := new(ClientMessage)
	if err := ctx.Bind(cm); err != nil {
		logger.Wrapf(err, "comment ctx bind err")
		helper.Fail(ctx, errors.New("参数缺失～"))
		return
	}
	sm := NewMessage(a, cm)
	res, err := db.ES.Index(
		consts.IndicesMessagesConst,
		bytes.NewReader(sm.Bytes()),
		db.ES.Index.WithDocumentID(sm.Id),
	)
	if err != nil {
		logger.Wrapf(err, "comment err")
		helper.ServerErr(ctx)
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		data, _ := ioutil.ReadAll(res.Body)
		logger.Errorf("comment es err: %s", string(data))
		return
	}
	helper.Success(ctx, sm)
	return
}

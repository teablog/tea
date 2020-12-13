package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/derror"
	"github.com/teablog/tea/internal/module/account"
	"github.com/teablog/tea/internal/validate"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	Channel        *channel
	ChannelMembers = &channelMembers{
		mu: &sync.RWMutex{},
		m:  make(map[string][]account.Account, 20),
	}
)

type channel struct {
	Creator   *account.Account  `json:"creator"`
	Members   []account.Account `json:"members"`
	Title     string            `json:"title"`
	CreatedAt time.Time         `json:"created_at"`
	Id        string            `json:"id"`
	Type      string            `json:"type"`
}

type channelWithMessage struct {
	channel
	Messages []ServerMessage `json:"messages"`
	Total    uint64          `json:"total"`
}

type channelResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source channel `json:"_source"`
			Id     string  `json:"_id"`
		} `json:"hits"`
	} `json:"hits"`
}

type serverMsgResponse struct {
	Hits struct {
		Total struct {
			Value uint64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source ServerMessage `json:"_source"`
			Id     string        `json:"_id"`
		} `json:"hits"`
	} `json:"hits"`
}

// 创建一个新的channel
func (*channel) Create(ctx *gin.Context, v *validate.ChannelCreateValidator) (c *channel, err error) {
	if a, ok := ctx.Get("account"); ok {
		acct := a.(*account.Account)
		type esResponse struct {
			Id string `json:"_id"`
		}
		// 创建channel
		var (
			buf   bytes.Buffer
			title = v.Title
		)
		// 标题设置
		m := account.NewAccount().Mget(v.Members)
		if len(*m) == 0 {
			return c, errors.Errorf("members not found")
		}
		for k, _ := range *m {
			// create_at 记录加入时间
			(*m)[k].CreateAt = time.Now()
		}
		if v.Type == consts.TypeChannelPublic && strings.TrimSpace(v.Title) == "" {
			name := make([]string, 0)
			for _, v := range *m {
				name = append(name, v.Name)
			}
			if len(name) > 20 {
				name = name[:20]
			}
			title = strings.Join(name, "、")
		}
		query := map[string]interface{}{
			"title":      title,
			"creator":    acct,
			"created_at": time.Now(),
			"type":       v.Type,
			"members":    *m,
		}
		if err = json.NewEncoder(&buf).Encode(query); err != nil {
			panic(errors.Wrap(err, "channel map json encode failed"))
		}
		res, err := db.ES.Index(
			consts.IndicesChannelConst,
			strings.NewReader(buf.String()),
			db.ES.Index.WithRefresh(`true`),
		)
		if err != nil {
			panic(errors.Wrap(err, "es create channel failed"))
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		if res.IsError() {
			panic(errors.Errorf("[%s] es create channel response: %s", res.StatusCode, body))
		}
		var r esResponse
		if err = json.Unmarshal(body, &r); err != nil {
			panic(errors.Wrapf(err, "response body json decode failed"))
		}
		c = &channel{
			Id:        r.Id,
			Members:   *m,
			Title:     title,
			Creator:   a.(*account.Account),
			CreatedAt: time.Now(),
			Type:      v.Type,
		}
		c.Title = c.getTitle(acct.Id)
	} else {
		panic(derror.Unauthorized{})
	}
	return
}

// 获取channel
func (*channel) Private(ctx *gin.Context, v *validate.ChannelCreateValidator) (c *channel, ok bool) {
	if a, ok := ctx.Get("account"); ok {
		acct := a.(*account.Account)
		query := fmt.Sprintf(`
{
  "query": {
    "bool": {
      "should": [
        {
          "bool": {
            "must": [
              {
                "term": {
                  "creator.id": "%s"
                }
              },
              {
                "term": {
                  "members.id": "%s"
                }
              },
              {
                "term": {
                  "type": "%s"
                }
              }
            ]
          }
        },
        {
          "bool": {
            "must": [
              {
                "term": {
                  "creator.id": "%s"
                }
              },
              {
                "term": {
                  "members.id": "%s"
                }
              },
              {
                "term": {
                  "type": "%s"
                }
              }
            ]
          }
        }
      ]
    }
  }
}`, acct.Id, v.Members[0], v.Type, v.Members[0], acct.Id, v.Type)
		res, err := db.ES.Search(
			db.ES.Search.WithIndex(consts.IndicesChannelConst),
			db.ES.Search.WithBody(strings.NewReader(query)),
		)
		if err != nil {
			panic(errors.Wrap(err, "channel exists search failed"))
		}
		defer res.Body.Close()
		respBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(errors.Wrapf(err, "[%d] es response body read failed", res.StatusCode))
		}
		if res.IsError() {
			panic(errors.Wrapf(err, "[%d] es response: %s", res.StatusCode, respBody))
		}
		var r channelResponse
		if err = json.Unmarshal(respBody, &r); err != nil {
			panic(errors.Wrapf(err, "channel exists es response: %s", respBody))
		}
		if r.Hits.Total.Value > 0 {
			c = &r.Hits.Hits[0].Source
			c.Id = r.Hits.Hits[0].Id
			// 私聊频道需要使用对方名称作为标题
			c.Title = c.getTitle(acct.Id)
			return c, true
		} else {
			return c, false
		}
	} else {
		panic(derror.Unauthorized{})
	}
}

func (*channel) Get(id string) (*channel, error) {
	type esResponse struct {
		Id     string  `json:"_id"`
		Source channel `json:"_source"`
	}
	res, err := db.ES.Get(
		consts.IndicesChannelConst,
		id,
	)
	if err != nil {
		panic(errors.Wrap(err, "es error"))
	}
	defer res.Body.Close()
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(errors.Errorf("[%d] es response body read error", res.StatusCode))
	}
	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			return nil, errors.Errorf("频道(%s)不存在", id)
		}
		panic(errors.Errorf("[%d] es response: %s", res.StatusCode, string(resp)))
	}
	var r esResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		return nil, errors.Wrapf(err, "es response: %s", string(resp))
	}
	s := &r.Source
	return s, nil
}

// 订阅channel
func (*channel) Subscribe(ctx *gin.Context) (*[]channel, error) {
	a, _ := ctx.Get("account")
	acct := a.(*account.Account)
	query := fmt.Sprintf(`
{
  "query": {
    "bool": {
      "should": [
        {"term": { "creator.id": "%s"}},
        {"term": { "members.id": "%s"}},
		{"term": { "_id": { "value": "douyacun", "boost": 10 }}}
      ]
    }
  }
}
`, acct.Id, acct.Id)
	res, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesChannelConst),
		db.ES.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		panic(errors.Wrap(err, "es channel index error"))
	}
	defer res.Body.Close()
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(errors.Wrap(err, "es response body read error"))
	}
	if res.IsError() {
		panic(errors.Errorf("[%d] es response: %s", res.StatusCode, resp))
	}
	var r channelResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		panic(errors.Wrapf(err, "json decode failed response: %s", resp))
	}
	var c []channel
	for _, v := range r.Hits.Hits {
		v.Source.Id = v.Id
		v.Source.Title = v.Source.getTitle(acct.Id)
		c = append(c, v.Source)
	}
	return &c, nil
}

func (c *channel) getTitle(accountId string) string {
	if c.Type == consts.TypeChannelPrivate {
		if c.Creator.Id == accountId {
			return c.Members[0].Name
		} else {
			return c.Creator.Name
		}
	} else {
		return c.Title
	}
}

func (c *channel) GetJoinTime(accountId string) time.Time {
	if c.Creator.Id == accountId {
		return c.Creator.CreateAt
	} else {
		for _, m := range c.Members {
			if m.Id == accountId {
				return m.CreateAt
			}
		}
	}
	return c.CreatedAt
}

type channelMembers struct {
	m  map[string][]account.Account
	mu *sync.RWMutex
}

func (m *channelMembers) Join(channelId string, accountId string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	members, ok := m.m[channelId]
	if !ok {
		c, err := Channel.Get(channelId)
		if err != nil {
			return err
		}
		members = c.Members
	}
	a, err := account.NewAccount().Get(accountId)
	if err != nil {
		return err
	}
	// create_at 记录加入时间
	(*a).CreateAt = time.Now()
	members = append(members, *a)
	m.m[channelId] = members
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(members); err != nil {
		panic(errors.Wrap(err, "members json encode error"))
	}
	res, err := db.ES.Update(
		consts.IndicesChannelConst,
		channelId,
		strings.NewReader(buf.String()),
	)
	if err != nil {
		panic(errors.Wrap(err, "channel members update error"))
	}
	defer res.Body.Close()
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(errors.Wrap(err, "channel members update response read error"))
	}
	if res.IsError() {
		panic(errors.Errorf("[%d] channel members update response: %s", res.StatusCode, resp))
	}
	return nil
}

func (m *channelMembers) Members(channelId string) (*[]account.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	members, ok := m.m[channelId]
	if !ok {
		ch, err := Channel.Get(channelId)
		if err != nil {
			return nil, err
		}
		members = ch.Members
		members = append(members, *ch.Creator)
		m.m[channelId] = members
	}
	return &members, nil
}

func (m *channelMembers) MembersIds(channelId string) ([]string, error) {
	mm, err := m.Members(channelId)
	if err != nil {
		return nil, err
	}
	var data []string
	for _, v := range *mm {
		data = append(data, v.Id)
	}
	return data, nil
}

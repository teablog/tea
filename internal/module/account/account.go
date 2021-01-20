package account

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/db"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Accouter interface {
	GetName() string
	GetEmail() string
	GetId() string
	GetAvatarUrl() string
	GetUrl() string
	Source() string
}

type Account struct {
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	AvatarUrl string    `json:"avatar_url"`
	Email     string    `json:"email"`
	CreateAt  time.Time `json:"create_at"`
	Ip        string    `json:"ip"`
}

func NewAccount() *Account {
	return &Account{}
}

func (a *Account) Create(ctx *gin.Context, i Accouter) (data *Account, err error) {
	id := helper.Md516([]byte(i.GetId() + i.Source()))
	ava := a.avatar(id, i.GetAvatarUrl())
	var buf bytes.Buffer
	data = &Account{
		Name:      i.GetName(),
		Source:    i.Source(),
		Id:        i.GetId(),
		Url:       i.GetUrl(),
		AvatarUrl: ava,
		Email:     i.GetEmail(),
		CreateAt:  time.Now(),
		Ip:        helper.RealIP(ctx.Request),
	}
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		panic(errors.Wrap(err, "Account create json encode failed"))
	}
	res, err := db.ES.Index(
		consts.IndicesAccountConst,
		strings.NewReader(buf.String()),
		db.ES.Index.WithDocumentID(id),
	)
	if err != nil {
		panic(errors.Wrap(err, "es Account index create failed"))
	}
	defer res.Body.Close()
	if res.IsError() {
		d, _ := ioutil.ReadAll(res.Body)
		panic(errors.Errorf("[%s] Account create failed, response: %s", res.StatusCode, d))
	}
	data.Id = id
	return
}

func (a *Account) avatar(id, url string) string {
	storageDir := config.Path.StorageDir()
	ext := path.Ext(url)
	dialer, err := proxy.SOCKS5("tcp", config.Proxy.Socks5(), nil, proxy.Direct)
	if err != nil {
		logger.Errorf("new proxy socks5 err: %s", err.Error())
		return url
	}
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: dialer.Dial,
		},
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Errorf("new http request err: %s", err.Error())
		return url
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("get %s err: %s", url, err.Error())
		return url
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		res, _ := ioutil.ReadAll(resp.Body)
		logger.Errorf("get url %s response %d %s", url, resp.StatusCode, res)
		return url
	}
	if ext == "" {
		switch resp.Header.Get("content-type") {
		case "image/jpeg":
			ext = ".jpeg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			ext = ".jpg"
		}
	}
	res := path.Join("/images/avatar", id+ext)
	storageFile := path.Join(path.Dir(storageDir), res)
	f, err := os.OpenFile(storageFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(errors.Wrapf(err, "open avatar file error"))
	}
	defer f.Close()
	if _, err = io.Copy(f, resp.Body); err != nil {
		panic(errors.Wrapf(err, "copy response to file error"))
	}
	return res
}

func (a *Account) EnableAccess() bool {
	res, err := db.ES.Exists(
		consts.IndicesAccountConst,
		a.Id,
	)
	if err != nil {
		panic(errors.Wrap(err, "account enable access ES exists failed"))
	}
	if res.IsError() {
		return false
	}
	return true
}

func Get(id string) (*Account, error) {
	type esResponse struct {
		Id     string  `json:"_id"`
		Source Account `json:"_source"`
	}
	res, err := db.ES.Index(
		consts.IndicesAccountConst,
		strings.NewReader(``),
		db.ES.Index.WithDocumentID(id),
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
			return nil, errors.Errorf("账户(%s)不存在", id)
		}
		panic(errors.Errorf("[%d] es response: %s", res.StatusCode, string(resp)))
	}
	var r esResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		panic(errors.Wrapf(err, "es response: %s", string(resp)))
	}
	s := &r.Source
	return s, nil
}

func Search(body string) (int, AcctSlice, error) {
	resp, err := db.ES.Search(
		db.ES.Search.WithIndex(consts.IndicesArticleCost),
		db.ES.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var eslist db.ESListResponse
	if err := json.NewDecoder(resp.Body).Decode(&eslist); err != nil {
		return 0, nil, err
	}
	type source struct {
		Source *Account `json:"_source"`
	}
	hits := make([]*source, 0)
	if err := json.Unmarshal(eslist.Hits.Hits, &hits); err != nil {
		return 0, nil, err
	}
	m := make(AcctSlice, 0)
	for _, v := range hits {
		m = append(m, v.Source)
	}
	return eslist.Hits.Total.Value, m, nil
}

package account

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/logger"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type _github struct {
	t struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	u struct {
		Id        int64  `json:"id"`
		Url       string `json:"html_url"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarUrl string `json:"avatar_url"`
	}
}

func NewGithub() *_github {
	return &_github{}
}

func (g *_github) TokenV2(code string) error {
	dialer, err := proxy.SOCKS5("tcp", config.Proxy.Socks5(), nil, proxy.Direct)
	if err != nil {
		return errors.Wrapf(err, "proxy socks5 err")
	}
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: dialer.Dial,
		},
		Timeout: 10 * time.Second,
	}
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(gin.H{
		"client_id":     config.Github.ClientId(),
		"client_secret": config.Github.ClientSecret(),
		"code":          code,
	}); err != nil {
		return errors.Wrapf(err, "github oauth access_token err")
	}
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", body)
	if err != nil {
		return errors.Wrapf(err, "github oauth access_token err")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	retries := 3
	for retries > 0 {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		retries--
	}
	if err != nil {
		panic(errors.Wrap(err, "client request error"))
	}
	if err != nil {
		return errors.Wrapf(err, "github oauth access_tokoen err")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrap(err, "response body read error"))
	}
	if resp.StatusCode > 299 {
		return errors.Errorf("github oauth access_token [%d] response: %s", resp.StatusCode, string(data))
	}
	if err := json.Unmarshal(data, &g.t); err != nil {
		panic(errors.Wrapf(err, "github oauth access_token err response: %s", data))
	}
	if g.t.AccessToken == "" {
		logger.Errorf("github oauth access_token [%d] response: %s", resp.StatusCode, string(data))
		return errors.Errorf("github授权登录失败")
	}
	return nil
}

func (g *_github) UserV2() error {
	authorization := bytes.NewBufferString(g.t.TokenType)
	authorization.WriteString(" ")
	authorization.WriteString(g.t.AccessToken)
	dialer, err := proxy.SOCKS5("tcp", config.Proxy.Socks5(), nil, proxy.Direct)
	if err != nil {
		return errors.Wrapf(err, "proxy socks5 err")
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
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return errors.Wrapf(err, "new http request err")
	}
	req.Header.Set("Authorization", authorization.String())
	var resp *http.Response
	retries := 3
	for retries > 0 {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		retries--
	}
	if err != nil {
		panic(errors.Wrap(err, "client request error"))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrapf(err, "read response body error"))
	}
	if resp.StatusCode > 299 {
		panic(errors.Errorf("[%d] response: %s", resp.StatusCode, string(data)))
	}
	if err = json.Unmarshal(data, &g.u); err != nil {
		panic(errors.Wrapf(err, "github user json decode error, response: %s", string(data)))
	}
	if g.u.Id == 0 {
		panic(errors.Wrapf(err, "github user access error, response: %s", string(data)))
	}
	return nil
}

func (g *_github) GetName() string {
	if strings.Trim(g.u.Name, " ") == "" {
		u, err := url.Parse(g.u.Url)
		if err != nil {
			return ""
		}
		match := strings.Split(strings.TrimLeft(u.Path, "/"), "/")
		if len(match) > 0 {
			return match[0]
		}
	}
	return g.u.Name
}

func (g *_github) GetId() string {
	return strconv.FormatInt(g.u.Id, 10)
}

func (g *_github) GetUrl() string {
	return g.u.Url
}

func (g *_github) GetEmail() string {
	return g.u.Email
}

func (g *_github) GetAvatarUrl() string {
	return g.u.AvatarUrl
}

func (g *_github) Source() string {
	return "github"
}

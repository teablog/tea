package account

import (
	"bytes"
	"dyc/internal/config"
	"dyc/internal/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

func (g *_github) Token(code string) (err error) {
	body := gin.H{
		"url":    "https://github.com/login/oauth/access_token",
		"method": "POST",
		"header": gin.H{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		"body": gin.H{
			"client_id": config.GetKey("github::client_id").String(),
			"client_secret": config.GetKey("github::client_secret").String(),
			"code": code,
		},
		"skip_verify": false,
		"timeout":     5 * time.Second,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		panic(errors.Wrapf(err, "body json encode error: %s", err.Error()))
	}
	req, _ := http.NewRequest("POST", config.GetKey("proxy::request").String(), &buf)
	client := http.Client{}
	retries := 3
	var resp *http.Response
	for retries > 0 {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		retries--
	}
	if err != nil {
		panic(errors.Wrap(err, "request github oauth failed"))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrap(err, "response body read error"))
	}
	if resp.StatusCode > 299 {
		return errors.Errorf("github oauth/access_token [%d] response: %s", resp.StatusCode, string(data))
	}
	if err := json.Unmarshal(data, &g.t); err != nil {
		panic(errors.Wrapf(err, "github oauth/access_token json decode error, response: %s", data))
	}
	if g.t.AccessToken == "" {
		logger.Errorf("oauth/access_token [%d] response: %s", resp.StatusCode, string(data))
		return errors.Errorf("github授权登录失败")
	}
	return
}

func (g *_github) User() (err error) {
	authorization := bytes.NewBufferString(g.t.TokenType)
	authorization.WriteString(" ")
	authorization.WriteString(g.t.AccessToken)
	query := gin.H{
		"url": "https://tea.github.com/user",
		"method": "GET",
		"header": gin.H{
			"Authorization": authorization.String(),
		},
		"skip_verify": false,
		"timeout": 5 * time.Second,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(errors.Wrapf(err, "user query json encode error"))
	}
	req, _ := http.NewRequest("POST", config.GetKey("proxy::request").String(), &buf)
	client := &http.Client{}
	retries := 3
	var resp *http.Response
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
	return
}

func (g *_github) GetName() string {
	if strings.Trim(g.u.Name, " ") == ""{
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

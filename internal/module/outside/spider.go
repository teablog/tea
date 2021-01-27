package outside

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/mail"
	"github.com/teablog/tea/internal/module/outside/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Spider 爬虫验证友情链接是否已经添加
func (row *Outside) Spider() {
	id := row.GenId()
	// 过滤黑名单
	if blacks, err := All(); err != nil {
		logger.Wrapf(err, "outside black es search err ")
		return
	} else {
		for _, v := range blacks {
			if strings.Contains(row.Url, v.Host) && v.Id != id {
				return
			}
		}
	}
	skill := fmt.Sprintf(`<p><a href="%s/article/%s">如何在Douyacun添加友情链接</a></p>`, config.Global.Host(), config.ES.FriendsLinkId())
	msg := mail.NewMessage()
	msg.SetTo(row.Email)
	msg.SetTitle("Douyacun 友情链接")
	up, err := url.Parse(row.Url)
	if err != nil {
		bd := "<p>友情链接添加失败了，您的网站似乎不是正常链接：</p>"
		bd += fmt.Sprintf("<p>网站地址: %s</p>", row.Url)
		bd += fmt.Sprintf("<p>%s</p>%s", err.Error(), skill)
		msg.SetBody(html.Common(bd))
		mail.Send(msg)
		return
	}
	body, err := fetchUrl(row.Url)
	if err != nil {
		logger.Debugf("fetch url err: %s", row.Url)
		bd := "<p>友情链接添加失败了，您的网站访问似乎出了点问题</p>"
		bd += fmt.Sprintf("<p>网站地址: %s</p>", row.Url)
		bd += fmt.Sprintf("<p>%s</p>%s", err.Error(), skill)
		msg.SetBody(html.Common(bd))
		mail.Send(msg)
		return
	}
	// 匹配 站点链接
	if err := matchHost(body, OutsideReg); err != nil {
		if err == ErrorNoMatch {
			bd := "<p>友情链接添加失败了，没有在您的网站中访问到Douyacun\n</p>"
			bd += fmt.Sprintf("<p>网站地址: %s</p>", row.Url)
			bd += skill
			msg.SetBody(html.Common(bd))
			mail.Send(msg)
			return
		}
		logger.Wrapf(err, "spider match host err: ")
		return
	}
	// 匹配 robots.txt
	if matchRobots(fmt.Sprintf("%s://%s/robots.txt", up.Scheme, up.Host)) {
		bd := "<p>友情链接添加失败了，没有在您的网站中访问到Douyacun</p>"
		bd += fmt.Sprintf("<p>网站地址: %s</p>", row.Url)
		bd += skill
		msg.SetBody(html.Common(bd))
		mail.Send(msg)
		return
	}
	if err := row.create(); err != nil {
		logger.Wrapf(err, "outside create err ")
		return
	}
	bd := "<p>友情链接添加成功</p>"
	bd += fmt.Sprintf(`<p>点此查看: <a href="%s">%s</a></p>`, config.Global.Host(), config.Global.Host())
	msg.SetBody(html.Common(bd))
	mail.Send(msg)
}

// match 匹配是否已经添加外链
func matchHost(body []byte, reg string) error {
	re, err := regexp.Compile(reg)
	if err != nil {
		return err
	}
	if !re.Match(body) {
		return ErrorNoMatch
	}
	ss := re.FindAllString(string(body), -1)
	ok := false
	for _, v := range ss {
		if !strings.Contains(v, "nofollow") {
			ok = true
		}
	}
	if !ok {
		return ErrorNoMatch
	}
	if !strings.Contains(string(body), "ouyacun") {
		return ErrorNoMatch
	}
	return nil
}

// fetchUrl 请求站点页面
func fetchUrl(url string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	clt := http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	resp, err := clt.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.Errorf("response status code %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return nil, errors.New("no content")
	}
	return data, nil
}

func matchRobots(url string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	clt := http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	resp, err := clt.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		data, _ := ioutil.ReadAll(resp.Body)
		return strings.Contains(string(data), "douyacun")
	}
	return false
}

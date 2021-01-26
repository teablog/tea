package outside

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/logger"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Spider 爬虫验证友情链接是否已经添加
func Spider(title, url, email string) {
	// todo 钉钉通知
	body, err := fetch(url)
	if err != nil {
		//	todo 请求异常报警邮件
		return
	}
	re, err := regexp.Compile(OutsideReg)
	if err != nil {
		logger.Errorf("regexp compile %s err", OutsideReg)
		return
	}
	if !re.Match(body) {
		// todo no match 报警邮件
		return
	}
	ss := re.FindAllString(string(body), -1)
	ok := false
	for _, v := range ss {
		if !strings.Contains(v, "nofollow") {
			ok = true
		}
	}
	if !ok {
		// todo nofollow报警邮件
		return
	}
}

func fetch(url string) ([]byte, error) {
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

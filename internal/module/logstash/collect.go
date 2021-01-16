package logstash

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/logger"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

type tries struct {
	BalancerLatency int64  `json:"balancer_latency"`
	Port            int64  `json:"port"`
	BalancerStart   int64  `json:"balancer_start"`
	Ip              string `json:"ip"`
}

type KongHttpLog struct {
	Latencies struct {
		Request json.Number `json:"request"`
		Kong    json.Number `json:"kong"`
		Proxy   json.Number `json:"proxy"`
	} `json:"latencies"`
	Request struct {
		Querystring map[string]string `json:"querystring"`
		Size        json.Number       `json:"size"`
		Uri         string            `json:"uri"`
		Url         string            `json:"url"`
		Headers     map[string]string `json:"headers"`
		Method      string            `json:"method"`
	} `json:"request"`
	ClientIp string  `json:"client_ip"`
	Tries    []tries `json:"tries"`
	Response struct {
		Headers map[string]string `json:"headers"`
		Status  json.Number       `json:"status"`
		Size    json.Number       `json:"size"`
	}
	Spider string `json:"spider"`
	Date   string `json:"date"`
}

func Accept(ctx *gin.Context) error {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	r := new(KongHttpLog)
	if err := json.Unmarshal(data, r); err != nil {
		return errors.Wrapf(err, "raw data: %s", data)
	}
	if sp, ok := match(r.Request.Headers["user-agent"]); ok {
		logger.Debugf("request body: %s", string(data))
		r.Spider = sp
		r.Date = time.Now().Format(consts.EsTimeFormat)
		if nData, err := json.Marshal(r); err != nil {
			return err
		} else {
			if err := ES.KongHttpLog(string(nData)); err != nil {
				return err
			}
		}
	}
	return nil
}

// 匹配爬虫
func match(spider string) (string, bool) {
	r, err := regexp.Compile(config.Logstash.Spider())
	if err != nil {
		return "", false
	}
	if !r.MatchString(spider) {
		return "", false
	}
	rs := r.FindStringSubmatch(spider)
	if len(rs) == 0 {
		return "", false
	}
	return rs[0], true
}

// 过滤日志
func filterUri(data *KongHttpLog) bool {
	for _, v := range config.Logstash.FilterUri() {
		if vv := strings.Trim(v, " "); vv != "" {
			if strings.HasPrefix(data.Request.Uri, vv) {
				return false
			}
		}
	}
	return true
}

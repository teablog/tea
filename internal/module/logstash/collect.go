package logstash

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/logger"
	"io/ioutil"
	"regexp"
	"strings"
)

type tries struct {
	BalancerLatency int64  `json:"balancer_latency"`
	Port            int64  `json:"port"`
	BalancerStart   int64  `json:"balancer_start"`
	Ip              string `json:"ip"`
}

type KongHttpLog struct {
	Latencies struct {
		Request int64 `json:"request"`
		Kong    int64 `json:"kong"`
		Proxy   int64 `json:"proxy"`
	} `json:"latencies"`
	Request struct {
		Querystring map[string]string `json:"querystring"`
		Size        int64             `json:"size"`
		Uri         string            `json:"uri"`
		Url         string            `json:"url"`
		Headers     map[string]string `json:"headers"`
		Method      string            `json:"method"`
	} `json:"request"`
	ClientIp string  `json:"client_ip"`
	Tries    []tries `json:"tries"`
	Response struct {
		Headers map[string]string `json:"headers"`
		Status  int64             `json:"status"`
		Size    int64             `json:"size"`
	}
	Spider string `json:"spider"`
}

func Accept(ctx *gin.Context) error {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	logger.Debugf("request body: %s", string(data))
	r := new(KongHttpLog)
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	if sp, ok := match(r.Request.Headers["user-agent"]); ok {
		r.Spider = sp
	}
	if !filterUri(r) {
		return nil
	}
	if nData, err := json.Marshal(r); err != nil {
		return err
	} else {
		if err := ES.KongHttpLog(string(nData)); err != nil {
			return err
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

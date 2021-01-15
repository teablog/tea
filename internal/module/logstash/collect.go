package logstash

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/config"
	"io/ioutil"
	"regexp"
)

type KongHttpLog struct {
	Request struct {
		Headers struct {
			UA string `json:"user-agent"`
		} `json:"headers"`
	} `json:"request"`
}

func Accept(ctx *gin.Context) error {
	//ctx.Request.Body
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	r := new(KongHttpLog)
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	if r.Request.Headers.UA == "" {
		return nil
	}
	if sp, ok := match(r.Request.Headers.UA); ok {
		if err := ES.KongHttpLog(string(data), sp); err != nil {
			return err
		}
	}
	return nil
}

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
	return rs[0], false
}

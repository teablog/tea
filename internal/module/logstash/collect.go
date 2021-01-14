package logstash

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

}

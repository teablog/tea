package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/logstash"
)

var Logstash *_logstash

type _logstash struct{}

func (*_logstash) Collect(ctx *gin.Context) {
	if err := logstash.Accept(ctx); err != nil {
		logger.Errorf("logstash.accept err: %s", err.Error())
	}
	helper.Success(ctx, "ok")
}

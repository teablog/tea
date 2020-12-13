package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/module/chat"
	"github.com/teablog/tea/internal/validate"
)

var Channel *_channel

type _channel struct{}

func (*_channel) Create(ctx *gin.Context) {
	var (
		v   *validate.ChannelCreateValidator
		err error
	)
	if err = ctx.ShouldBindJSON(&v); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if err = validate.DoValidate(v); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if v.Type == consts.TypeChannelPrivate {
		if c, ok := chat.Channel.Private(ctx, v); ok {
			helper.Success(ctx, c)
			return
		}
	}
	c, err := chat.Channel.Create(ctx, v)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, c)
}

package controllers

import (
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/account"
	"github.com/teablog/tea/internal/module/chat"
	"github.com/teablog/tea/internal/validate"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
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

func (*_channel) subscribe(ctx *gin.Context) {
	q := ctx.DefaultQuery("q", "")
	r := make(map[string]time.Time, 0)
	if q != "" {
		if err := json.Unmarshal([]byte(q), &r); err != nil {
			helper.Fail(ctx, err)
			return
		}
	}
	data, err := chat.Channel.SubscribeWithMsg(ctx, &r)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, data)
	return
}

func (*_channel) Messages(ctx *gin.Context) {
	var vld validate.ChannelMessagesValidator
	if err := ctx.ShouldBindQuery(&vld); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if err := validate.DoValidate(vld); err != nil {
		helper.Fail(ctx, err)
		return
	}
	sec := int64(vld.Before/1000)
	nsec := int64(vld.Before % 1000) * 1000000
	before := time.Unix(sec, nsec)
	logger.Debug(before.Zone())

	channel, err := chat.Channel.Get(vld.ChannelId)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	acct, _ := ctx.Get("account")
	joinDate := channel.GetJoinTime(acct.(*account.Account).Id)
	total, messages := chat.Channel.Messages(vld.ChannelId, joinDate, before)
	helper.Success(ctx, gin.H{"total": total, "messages": messages})
	return
}

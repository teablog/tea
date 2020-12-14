package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/module/chat"
	"strings"
)

var WS _ws

type _ws struct{}

func (*_ws) Join(ctx *gin.Context, hub *chat.Hub) {
	articleId := ctx.Query("article_id")
	if len(strings.Trim(articleId, " \r\n")) == 0 {
		helper.Fail(ctx, errors.New("参数缺失: article_id"))
		return
	}
	chat.ServeWs(ctx, hub, articleId)
}

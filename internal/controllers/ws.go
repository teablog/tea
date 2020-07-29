package controllers

import (
	"github.com/teablog/tea/internal/module/chat"
	"github.com/gin-gonic/gin"
)

var WS _ws

type _ws struct{}

func (*_ws) Join(ctx *gin.Context, hub *chat.Hub) {
	chat.ServeWs(ctx, hub)
}
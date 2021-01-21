package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/module/ws"
)

var WS _ws

type _ws struct{}

func (*_ws) Join(ctx *gin.Context, hub *ws.Hub) {

}

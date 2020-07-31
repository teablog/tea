package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/module/setting"
)

var Help *help

type help struct{}

func (h *help) Token(ctx *gin.Context) {
	raw := setting.Token.Get(ctx)
	helper.Success(ctx, raw)
}

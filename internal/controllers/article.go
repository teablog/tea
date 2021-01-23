package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/module/message"
	"github.com/teablog/tea/internal/validate"
)

var (
	Article _article
)

type _article struct{}


func (*_article) Messages(ctx *gin.Context) {
	var vld validate.MessagesValidator
	if err := ctx.ShouldBindQuery(&vld); err != nil {
		helper.Fail(ctx, err)
		return
	}
	if err := validate.DoValidate(vld); err != nil {
		helper.Fail(ctx, err)
		return
	}
	total, data, err := message.FindMessages(vld)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, gin.H{"total": total, "list": data})
}

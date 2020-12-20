package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/article"
	"github.com/teablog/tea/internal/module/chat"
	"github.com/teablog/tea/internal/validate"
	"net/http"
	"strconv"
	"strings"
)

var (
	Article _article
)

type _article struct{}

func (*_article) List(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	total, data, err := article.Post.List(c, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "服务器出错了!")
		return
	}
	helper.Success(c, gin.H{"total": total, "data": data})
	return
}

func (*_article) Labels(c *gin.Context) {
	// 关键字数量
	size := 30
	count := strings.TrimSpace(c.Param("count"))
	if count != "" {
		n, err := strconv.Atoi(count)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "非法请求"})
		}
		size = n
	}
	labels, _ := article.Label.List(size)
	helper.Success(c, labels)
	return
}

func (*_article) View(c *gin.Context) {
	id := c.Param("id")
	data, err := article.Post.View(c, id)
	if err != nil {
		logger.Errorf("%s", err)
		helper.Fail(c, errors.New("服务器出错了"))
		return
	}
	helper.Success(c, gin.H{"data": data})
}

func (*_article) Search(c *gin.Context) {
	q := c.Query("q")
	if len(q) == 0 {
		helper.Fail(c, errors.New("请指定查询内容"))
		return
	}
	total, data, err := article.Search.List(q)
	if err != nil {
		logger.Errorf("文章搜索错误: %s", err)
		helper.Fail(c, errors.New("文章搜索出错了"))
		return
	}
	helper.Success(c, gin.H{"total": total, "data": data})
	return
}

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

	total, data, err := chat.Message.FindMessages(vld)
	if err != nil {
		helper.Fail(ctx, err)
		return
	}
	helper.Success(ctx, gin.H{"total": total, "list": data})
}

func (*_article) Comment(hub *chat.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var message chat.ClientMessage
		if err := ctx.Bind(&message); err != nil {
			helper.Fail(ctx, err)
			return
		}
		if err := chat.Message.SendMessage(ctx, hub, message); err != nil {
			helper.Fail(ctx, err)
			return
		}
		helper.Success(ctx, "success")
	}
}

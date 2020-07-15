package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/logger"
	"dyc/internal/module/article"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	Topic _topic
)

type _topic struct{}

func (*_topic) List(c *gin.Context) {
	topic := c.Param("topic")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	total, data, err := article.Topics.List(topic, page)
	if err != nil {
		logger.Errorf("首页文章列表错误: %s", err)
		helper.Fail(c, errors.New("服务器出错了"))
		return
	}
	helper.Success(c, gin.H{"total": total, "data": data})
	return
}

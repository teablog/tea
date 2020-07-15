package controllers

import (
	"dyc/internal/helper"
	"dyc/internal/logger"
	"dyc/internal/module/article"
	"errors"
	"github.com/gin-gonic/gin"
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

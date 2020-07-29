package controllers

import (
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/subscribe"
	"errors"
	"github.com/gin-gonic/gin"
	"regexp"
)

var Subscribe _subscribe

type _subscribe struct{}

type formEmail struct {
	Email string `json:"email"`
}

func (*_subscribe) Create(c *gin.Context) {
	var (
		err error
		e   formEmail
	)
	if err = c.ShouldBind(&e); err != nil {
		helper.Fail(c, err)
		return
	}
	logger.Debugf("%v", e)
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", e.Email); !m {
		helper.Fail(c, errors.New("邮件格式错误"))
		return
	}
	if err = subscribe.Email.Store(e.Email); err != nil {
		helper.Fail(c, errors.New("存储失败"))
		return
	}
	helper.Success(c, "success")
}

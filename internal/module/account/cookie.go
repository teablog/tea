package account

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/consts"
	"github.com/teablog/tea/internal/helper"
)

type Cookie struct {
	*Account
	Md5 string `json:"md5"`
}

func (a *Account) SetCookie(ctx *gin.Context) {
	var (
		c   Cookie
		err error
	)
	c.Account = a
	data, err := json.Marshal(a)
	if err != nil {
		panic(errors.Wrap(err, "set cookie json encode failed"))
	}
	c.Md5 = helper.Md532(data)
	cookie, err := json.Marshal(c)
	if err != nil {
		panic(errors.Wrap(err, "set cookie json encode failed"))
	}
	ctx.SetCookie(consts.CookieName, string(cookie), 31536000, "/", "."+config.Global.Domain(), false, false)
}

func (a *Account) ExpireCookie(ctx *gin.Context) {
	ctx.SetCookie(consts.CookieName, "", -1, "/", "."+config.Global.Domain(), false, false)
}

func (c *Cookie) VerifyCookie() bool {
	// 验证cookie完整性
	a, err := json.Marshal(c.Account)
	if err != nil {
		return false
	}
	if c.Md5 != helper.Md532(a) {
		return false
	}
	return true
}

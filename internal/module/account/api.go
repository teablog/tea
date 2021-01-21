package account

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"strings"
)

func NameExists(ctx *gin.Context) {
	name := ctx.Query("name")
	name = strings.Trim(name, " \r\n")
	if name == "" {
		helper.Fail(ctx, errors.New("起个像样点的名字吧～"))
	}
	query := fmt.Sprintf(`{"query":{"term":{"name":"%s"}}}`, name)
	c, _, err := Search(query)
	if err != nil {
		logger.Wrapf(err, "name exists search err")
		helper.Fail(ctx, errors.New("服务器出错了～"))
		return
	}
	helper.Success(ctx, c == 0)
	return
}

func Register(ctx *gin.Context) {
	acct := NewAccount()
	if err := ctx.Bind(acct); err != nil {
		logger.Errorf("参数异常: %s", err.Error())
		helper.Fail(ctx, errors.New("参数异常~"))
		return
	}
	u, err := acct.Create(ctx)
	if err != nil {
		logger.Wrapf(err, "account register err")
		helper.Fail(ctx, errors.New("服务异常～"))
		return
	}
	u.SetCookie(ctx)
	helper.Success(ctx, u)
	return
}

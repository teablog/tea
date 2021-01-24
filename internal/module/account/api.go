package account

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/helper"
	"github.com/teablog/tea/internal/logger"
	"net/http"
	"strings"
)

func Route(api *gin.RouterGroup) {
	a := api.Group("account")
	{
		a.GET("/name/exists", NameExists)
		a.POST("/register", Register)
	}
}

func NameExists(ctx *gin.Context) {
	name := strings.Trim(ctx.Query("name"), " \r\n")
	email := strings.Trim(ctx.Query("email"), " \r\n")
	if name == "" {
		resp := helper.Response{
			Code:    403,
			Message: "起个名字可真难～",
			Data:    "name",
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if email == "" {
		resp := helper.Response{
			Code:    0,
			Message: "邮箱还是要写一下的",
			Data:    "email",
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}
	query := fmt.Sprintf(`{"query":{"term":{"name":"%s"}}}`, name)
	c, _, err := Search(query)
	if err != nil {
		logger.Wrapf(err, "name exists search err")
		helper.Fail(ctx, errors.New("服务器出错了～"))
		return
	}
	helper.Success(ctx, c != 0)
	return
}

func Register(ctx *gin.Context) {
	acct := NewAccount()
	if err := ctx.Bind(acct); err != nil {
		logger.Errorf("参数异常: %s", err.Error())
		helper.Fail(ctx, errors.New("参数异常~"))
		return
	}
	acct.Name = strings.Trim(acct.Name, " \n")
	acct.Email = strings.Trim(acct.Email, " \n")
	if acct.Name == "" {
		resp := helper.Response{
			Code:    403,
			Message: "起个名字可真难～",
			Data:    "name",
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if acct.Email == "" {
		resp := helper.Response{
			Code:    403,
			Message: "邮箱还是要写一下的",
			Data:    "email",
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}
	_, accts, err := Search(fmt.Sprintf(`{"query":{"term":{"name":"%s"}}}`, acct.Name))
	if err != nil {
		logger.Wrapf(err, "account register search err")
		helper.ServerErr(ctx)
		return
	}
	if accts.NameRepeat(acct.Name, acct.Email) {
		resp := helper.Response{
			Code:    0,
			Message: fmt.Sprintf("%s已经被别的村民占用咧～", acct.Name),
			Data:    "name",
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}
	// 已经注册过了直接返回
	if res, ok := accts.Acct(acct.Name, acct.Email); ok {
		helper.Success(ctx, res)
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

package middleware

import (
	"dyc/internal/derror"
	"dyc/internal/logger"
	"dyc/internal/module/account"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/url"
)

func LoginCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		abort := func() {
			q := url.Values{}
			q.Add("redirect_uri", "/chat")
			account.NewAccount().ExpireCookie(ctx)
			panic(derror.Unauthorized{})
		}
		cookieStr, err := ctx.Cookie("douyacun")
		logger.Debugf("cookie: %s", cookieStr)
		if err != nil || cookieStr == "" {
			abort()
			return
		}
		// 验证cookie合法性
		var cookie account.Cookie
		if err = json.Unmarshal([]byte(cookieStr), &cookie); err != nil {
			abort()
			return
		}
		if !cookie.VerifyCookie() || !cookie.Account.EnableAccess() {
			abort()
			return
		}
		ctx.Set("account", cookie.Account)
		ctx.Next()
	}
}

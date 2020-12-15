package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/teablog/tea/internal/logger"
	"github.com/teablog/tea/internal/module/account"
	"net/http"
)

func LoginCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieStr, err := ctx.Cookie("douyacun")
		logger.Debugf("cookie: %s", cookieStr)
		if err != nil || cookieStr == "" {
			ctx.XML(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		// 验证cookie合法性
		var cookie account.Cookie
		if err = json.Unmarshal([]byte(cookieStr), &cookie); err != nil {
			account.NewAccount().ExpireCookie(ctx)
			ctx.XML(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		if !cookie.VerifyCookie() || !cookie.Account.EnableAccess() {
			ctx.XML(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			ctx.Abort()
			return
		}
		ctx.Set("account", cookie.Account)
		ctx.Next()
	}
}
